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
	"github.com/itsubaki/ghz/pkg/releases"
	"github.com/itsubaki/ghz/pkg/tags"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	projectID := dataset.ProjectID

	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID)

	token, err := NextToken(ctx, projectID, dsn)
	if err != nil {
		log.Error("next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("next token=%v", token)

	t, err := tags.Fetch(ctx,
		&tags.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
		},
	)
	if err != nil {
		log.Error("fetch tags: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("tags=%v", t)

	tags := make(map[string]*github.RepositoryTag)
	for i := range t {
		tags[t[i].GetName()] = t[i]
	}

	if _, err := releases.Fetch(ctx,
		&releases.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			LastID:     token,
		},
		func(list []*github.RepositoryRelease) error {
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

			if err := dataset.Insert(ctx, dsn, dataset.ReleasesMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		},
	); err != nil {
		log.Error("fetch: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func NextToken(ctx context.Context, projectID, dsn string) (int64, error) {
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
