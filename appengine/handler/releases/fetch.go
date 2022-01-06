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
	"github.com/itsubaki/ghz/pkg/releases"
	"github.com/itsubaki/ghz/pkg/tags"
)

type Response struct {
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	token, err := NextToken(ctx, id, dsn)
	if err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

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
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("fetch tags: %v", err),
		})
		return
	}

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
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("fetch: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Path: c.Request.URL.Path,
	})
}

func NextToken(ctx context.Context, projectID, datasetName string) (int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.ReleasesMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v`", table)

	var id int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		id = values[0].(int64)
	}); err != nil {
		return -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return id, nil
}
