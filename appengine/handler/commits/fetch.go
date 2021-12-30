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
	"github.com/itsubaki/ghz/appengine/dataset/view"
	"github.com/itsubaki/ghz/pkg/commits"
)

type Response struct {
	Path      string `json:"path"`
	NextToken string `json:"next_token"`
	Message   string `json:"message,omitempty"`
}

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.CommitsMeta,
		dataset.IncidentsMeta,
		dataset.PullReqsMeta,
		view.IncidentsCommitsMeta(id, dsn),
		view.IncidentsPullReqsMeta(id, dsn),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	token, err := NextToken(ctx, id, dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

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
		}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:      c.Request.URL.Path,
			NextToken: token,
			Message:   fmt.Sprintf("fetch: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Path:      c.Request.URL.Path,
		NextToken: token,
	})
}

func NextToken(ctx context.Context, projectID, datasetName string) (string, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.CommitsMeta.Name)
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
