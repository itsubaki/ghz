package releases

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/pkg/releases"
	"github.com/itsubaki/ghz/pkg/tags"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, dsn, []bigquery.TableMetadata{
		dataset.ReleasesMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := NextToken(ctx, id, dsn)
	if err != nil {
		log.Printf("get lastSHA: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, token)

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
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
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
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
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
