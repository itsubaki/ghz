package commits

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/commits"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.CommitsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := NextToken(ctx, datasetName)
	if err != nil {
		log.Printf("get lastSHA: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, token)

	if _, err := commits.Fetch(ctx,
		&commits.ListInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			LastSHA:    token,
		},
		func(list []*github.RepositoryCommit) error {
			items := make([]interface{}, 0)
			for _, r := range list {
				items = append(items, dataset.Commits{
					Owner:      owner,
					Repository: repository,
					SHA:        r.GetSHA(),
					Login:      r.Commit.Author.GetName(),
					Date:       r.Commit.Author.GetDate(),
					Message:    strings.ReplaceAll(r.Commit.GetMessage(), "\n", " "),
				})
			}

			if err := dataset.Insert(ctx, datasetName, dataset.CommitsTableMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		}); err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func NextToken(ctx context.Context, datasetName string) (string, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return "", fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.CommitsTableMeta.Name)
	query := fmt.Sprintf("select sha from `%v` where date = (select max(date) from `%v` limit 1)", table, table)

	var sha string
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		sha = values[0].(string)
	}); err != nil {
		return "", fmt.Errorf("query(%v): %v", query, err)
	}

	return sha, nil
}
