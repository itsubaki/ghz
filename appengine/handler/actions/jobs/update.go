package jobs

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
	"github.com/itsubaki/ghz/appengine/dataset/view"
	"github.com/itsubaki/ghz/pkg/actions/jobs"
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
		dataset.WorkflowRunsMeta,
		dataset.WorkflowJobsMeta,
		view.WorkflowJobsMeta(id, dsn),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	list, err := ListJobs(ctx, id, dsn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateResponse{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("list jobs: %v", err),
		})
		return
	}

	for _, j := range list {
		job, err := jobs.Get(ctx, &jobs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			JobID:      j.JobID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, UpdateResponse{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("get job(%v): %v", j.JobID, err),
			})
			continue
		}

		if err := UpdateJob(ctx, id, dsn, job); err != nil {
			log.Printf("%#v", UpdateResponse{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("update job(%v): %v", j.JobID, err),
			})
			continue
		}
	}

	c.JSON(http.StatusOK, UpdateResponse{
		Path: c.Request.URL.Path,
	})
}

func ListJobs(ctx context.Context, projectID, datasetName string) ([]dataset.WorkflowJob, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowJobsMeta.Name)
	query := fmt.Sprintf("select job_id from `%v` where status != \"completed\"", table)

	out := make([]dataset.WorkflowJob, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		out = append(out, dataset.WorkflowJob{
			JobID: values[0].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}

func UpdateJob(ctx context.Context, projectID, datasetName string, j *github.WorkflowJob) error {
	if j.GetStatus() != "completed" {
		return nil
	}

	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowJobsMeta.Name)
	query := fmt.Sprintf("update %v set status = \"%v\", conclusion = \"%v\", completed_at = \"%v\" where job_id = %v",
		table,
		j.GetStatus(),
		j.GetConclusion(),
		j.GetCompletedAt().Format("2006-01-02 15:04:05 UTC"),
		j.GetID(),
	)

	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		return
	}); err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	return nil
}
