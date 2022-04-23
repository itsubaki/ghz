package releases

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/appengine/tracer"
	"github.com/itsubaki/ghz/pkg/releases"
	"github.com/itsubaki/ghz/pkg/tags"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var tra = otel.Tracer("handler/releases/fetch")

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")
	traceTrue := c.GetBool("trace_true")

	projectID := dataset.ProjectID
	dsn := dataset.Name(owner, repository)

	log := logger.New(projectID, traceID).NewReport(ctx, c.Request)
	log.DebugWith(spanID, "trace_id=%v, span_id=%v, trace_true=%v", traceID, spanID, traceTrue)

	parent, err := tracer.NewContext(ctx, traceID, spanID, traceTrue)
	if err != nil {
		log.ErrorReport("new context: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var token int64
	if err := tracer.Span(tra, parent, "get next token", func(c context.Context, s trace.Span) error {
		token, err = GetNextToken(c, projectID, dsn)
		if err != nil {
			return err
		}

		log.DebugWith(s.SpanContext().SpanID().String(), "next token=%v", token)
		return nil
	}); err != nil {
		log.ErrorReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var rtags map[string]*github.RepositoryTag
	if err := tracer.Span(tra, parent, "fetch tags", func(c context.Context, s trace.Span) error {
		t, err := tags.Fetch(c,
			&tags.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
			},
		)
		if err != nil {
			return err
		}

		rtags = make(map[string]*github.RepositoryTag)
		for i := range t {
			rtags[t[i].GetName()] = t[i]
		}

		log.DebugWith(s.SpanContext().SpanID().String(), "fetched len(tags)=%v", len(rtags))
		return nil
	}); err != nil {
		log.ErrorReport("fetch tags: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := tracer.Span(tra, parent, "fetch releases", func(c context.Context, s trace.Span) error {
		if _, err := releases.Fetch(c,
			&releases.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
				LastID:     token,
			},
			func(list []*github.RepositoryRelease) error {
				return tracer.Span(tra, c, "insert items", func(c context.Context, s trace.Span) error {
					items := make([]interface{}, 0)
					for _, r := range list {
						items = append(items, dataset.Release{
							Owner:           owner,
							Repository:      repository,
							ID:              r.GetID(),
							TagName:         r.GetTagName(),
							TagSHA:          rtags[r.GetTagName()].GetCommit().GetSHA(),
							Login:           r.GetAuthor().GetLogin(),
							TargetCommitish: r.GetTargetCommitish(),
							Name:            r.GetName(),
							CreatedAt:       r.GetCreatedAt().Time,
							PublishedAt:     r.GetPublishedAt().Time,
						})
					}

					if err := dataset.Insert(c, dsn, dataset.ReleasesMeta.Name, items); err != nil {
						return fmt.Errorf("insert items: %v", err)
					}

					log.DebugWith(s.SpanContext().SpanID().String(), "inserted. len(items)=%v", len(items))
					return nil
				})
			},
		); err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.ErrorReport("fetch releases: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func GetNextToken(ctx context.Context, projectID, dsn string) (int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.ReleasesMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v`", table)

	var rid int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		rid = values[0].(int64)
	}); err != nil {
		return -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return rid, nil
}
