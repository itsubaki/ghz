package runs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/pkg/actions/runs"
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

	token, _, err := NextToken(ctx, id, dsn)
	if err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

	if _, err := runs.Fetch(ctx,
		&runs.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			LastID:     token,
		},
		func(list []*github.WorkflowRun) error {
			items := make([]interface{}, 0)
			for _, r := range list {
				items = append(items, dataset.WorkflowRun{
					Owner:        owner,
					Repository:   repository,
					WorkflowID:   r.GetWorkflowID(),
					WorkflowName: r.GetName(),
					RunID:        r.GetID(),
					RunNumber:    int64(r.GetRunNumber()),
					Status:       r.GetStatus(),
					Conclusion:   r.GetConclusion(),
					CreatedAt:    r.GetCreatedAt().Time,
					UpdatedAt:    r.GetUpdatedAt().Time,
					HeadSHA:      r.GetHeadSHA(),
				})
			}

			if err := dataset.Insert(ctx, dsn, dataset.WorkflowRunsMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		}); err != nil {
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

func NextToken(ctx context.Context, id, dsn string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", id, dsn, dataset.WorkflowRunsMeta.Name)
	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v` limit 1", table)

	var rid, num int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		rid = values[0].(int64)
		num = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return rid, num, nil
}
