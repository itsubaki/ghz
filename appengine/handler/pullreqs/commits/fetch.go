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
	Path      string `json:"path"`
	NextToken int64  `json:"next_token"`
	Message   string `json:"message,omitempty"`
}

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.PullReqsMeta,
		dataset.PullReqCommitsMeta,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	token, _, err := NextToken(ctx, id, dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

	prs, err := ListPullReqs(ctx, id, dsn, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:      c.Request.URL.Path,
			NextToken: token,
			Message:   fmt.Sprintf("list pullreqs: %v", err),
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
			c.JSON(http.StatusInternalServerError, Response{
				Path:      c.Request.URL.Path,
				NextToken: token,
				Message:   fmt.Sprintf("fetch: %v", err),
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
			c.JSON(http.StatusInternalServerError, Response{
				Path:      c.Request.URL.Path,
				NextToken: token,
				Message:   fmt.Sprintf("insert items: %v", err),
			})
			return
		}
	}

	c.JSON(http.StatusOK, Response{
		Path:      c.Request.URL.Path,
		NextToken: token,
	})
}

type PullReq struct {
	ID     int64
	Number int64
}

func ListPullReqs(ctx context.Context, projectID, datasetName string, nextToken int64) ([]PullReq, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.PullReqsMeta.Name)
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

func NextToken(ctx context.Context, projectID, datasetName string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.PullReqCommitsMeta.Name)
	query := fmt.Sprintf("select max(id), max(number) from `%v` limit 1", table)

	var id, num int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		id = values[0].(int64)
		num = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return id, num, nil
}
