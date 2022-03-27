package commits

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/pkg/pullreqs/commits"
)

type Response struct {
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	projectID := c.GetString("project_id")
	dsn := dataset.Name(owner, repository)

	token, _, err := NextToken(ctx, projectID, dsn)
	if err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

	prs, err := ListPullReqs(ctx, projectID, dsn, token)
	if err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("list pullreqs: %v", err),
		})
		return
	}

	for _, p := range prs {
		list, err := commits.Fetch(ctx,
			&commits.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
			},
			int(p.Number),
		)
		if err != nil {
			c.Error(err).SetMeta(Response{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("fetch: %v", err),
			})
			return
		}

		items := make([]interface{}, 0)
		for _, r := range list {
			message := regexpnl.ReplaceAllString(r.Commit.GetMessage(), " ")
			if len(message) > 64 {
				message = message[0:64]
			}

			items = append(items, dataset.PullReqCommits{
				Owner:      owner,
				Repository: repository,
				ID:         p.ID,
				Number:     p.Number,
				SHA:        r.GetSHA(),
				Login:      r.Commit.Author.GetName(),
				Date:       r.Commit.Author.GetDate(),
				Message:    message,
			})
		}

		if err := dataset.Insert(ctx, dsn, dataset.PullReqCommitsMeta.Name, items); err != nil {
			c.Error(err).SetMeta(Response{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("insert items: %v", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Path: c.Request.URL.Path,
	})
}

type PullReq struct {
	ID     int64
	Number int64
}

func ListPullReqs(ctx context.Context, projectID, dsn string, nextToken int64) ([]PullReq, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where id > %v", table, nextToken)

	prs := make([]PullReq, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		prs = append(prs, PullReq{
			ID:     values[0].(int64),
			Number: values[1].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return prs, nil
}

func NextToken(ctx context.Context, projectID, dsn string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqCommitsMeta.Name)
	query := fmt.Sprintf("select max(id), max(number) from `%v` limit 1", table)

	var pid, num int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		pid = values[0].(int64)
		num = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return pid, num, nil
}
