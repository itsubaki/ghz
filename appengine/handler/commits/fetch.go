package commits

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/commits"
)

var (
	projectID = dataset.ProjectID
	logf      = logger.Factory
	regexpnl  = regexp.MustCompile(`\r\n|\r|\n`)
)

func Fetch(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logf.New(traceID, c.Request)

	token, err := GetNextToken(ctx, projectID, dsn)
	if err != nil {
		log.ErrorReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("next token=%v", token)

	if _, err := commits.Fetch(ctx,
		&commits.FetchInput{
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
				message := regexpnl.ReplaceAllString(r.Commit.GetMessage(), " ")
				if len(message) > 64 {
					message = message[0:64]
				}

				items = append(items, dataset.Commit{
					Owner:      owner,
					Repository: repository,
					SHA:        r.GetSHA(),
					Login:      r.Commit.Author.GetName(),
					Date:       r.Commit.Author.GetDate(),
					Message:    message,
				})
			}

			if err := dataset.Insert(ctx, dsn, dataset.CommitsMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		},
	); err != nil {
		log.ErrorReport("fetch tags: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func GetNextToken(ctx context.Context, projectID, dsn string) (string, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.CommitsMeta.Name)
	query := fmt.Sprintf("select sha from `%v` where date = (select max(date) from `%v` limit 1)", table, table)

	var sha string
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
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
