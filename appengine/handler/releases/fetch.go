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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var (
	projectID = dataset.ProjectID
	logf      = logger.Factory
	tra       = otel.Tracer("handler/releases/fetch")
)

func Fetch(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")
	traceTrue := c.GetBool("trace_true")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logf.New(traceID, c.Request)
	log.SpanOf(spanID).Debug("trace=%v", traceTrue)

	parent, err := tracer.NewContext(ctx, traceID, spanID, traceTrue)
	if err != nil {
		log.ErrorReport("new context: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	token, err := func() (int64, error) {
		c, s := tra.Start(parent, "get next token",
			trace.WithAttributes(attribute.String("dataset_name", dsn)),
		)
		defer s.End()

		return GetNextToken(c, projectID, dsn)
	}()
	if err != nil {
		log.ErrorReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.SpanOf(spanID).Debug("next token=%v", token)

	t, err := func() ([]*github.RepositoryTag, error) {
		c, s := tra.Start(parent, "fetch tags")
		defer s.End()

		return tags.Fetch(c,
			&tags.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
			},
		)
	}()
	if err != nil {
		log.ErrorReport("fetch tags: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tags := make(map[string]*github.RepositoryTag)
	for i := range t {
		tags[t[i].GetName()] = t[i]
	}
	log.SpanOf(spanID).Debug("len(tags)=%v", len(tags))

	if _, err := func() ([]*github.RepositoryRelease, error) {
		c, s := tra.Start(parent, "fetch releases")
		defer s.End()

		return releases.Fetch(c,
			&releases.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
				LastID:     token,
			},
			func(list []*github.RepositoryRelease) error {
				c, s := tra.Start(c, "insert items")
				defer s.End()

				items := make([]interface{}, 0)
				for _, r := range list {
					items = append(items, dataset.Release{
						Owner:           owner,
						Repository:      repository,
						ID:              r.GetID(),
						TagName:         r.GetTagName(),
						TagSHA:          tags[r.GetTagName()].GetCommit().GetSHA(),
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

				log.Span(s).Debug("inserted. len(items)=%v", len(items))
				return nil
			})
	}(); err != nil {
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
