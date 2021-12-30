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
	"github.com/itsubaki/ghz/appengine/dataset/view"
	"github.com/itsubaki/ghz/pkg/actions/runs"
)

type UpdateResponse struct {
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

func Update(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.CommitsMeta,
		dataset.EventsPushMeta,
		dataset.PullReqsMeta,
		dataset.PullReqCommitsMeta,
		dataset.WorkflowRunsMeta,
		view.WorkflowRunsMeta(id, dsn),
		view.LeadTimeWorkflowsMeta(id, dsn),
		view.LeadTimeCommitsMeta(id, dsn),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	list, err := ListRuns(ctx, id, dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateResponse{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("list jobs: %v", err),
		})
		return
	}

	for _, r := range list {
		run, err := runs.Get(ctx, &runs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			RunID:      r.RunID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, UpdateResponse{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("get run(%v): %v", r.RunID, err),
			})
			return
		}

		if err := UpdateRun(ctx, id, dsn, run); err != nil {
			c.JSON(http.StatusInternalServerError, UpdateResponse{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("update run(%v): %v", r.RunID, err),
			})
			return
		}
	}
}

func ListRuns(ctx context.Context, projectID, datasetName string) ([]dataset.WorkflowRun, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowRunsMeta.Name)
	query := fmt.Sprintf("select run_id from `%v` where status != \"completed\"", table)

	out := make([]dataset.WorkflowRun, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		out = append(out, dataset.WorkflowRun{
			RunID: values[0].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func UpdateRun(ctx context.Context, projectID, datasetName string, r *github.WorkflowRun) error {
	if r.GetStatus() != "completed" {
		return nil
	}

	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowRunsMeta.Name)
	query := fmt.Sprintf("update %v set status = \"%v\", conclusion = \"%v\", updated_at = \"%v\" where run_id = %v",
		table,
		r.GetStatus(),
		r.GetConclusion(),
		r.GetUpdatedAt().Format("2006-01-02 15:04:05 UTC"),
		r.GetID(),
	)

	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		return
	}); err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	return nil
}
