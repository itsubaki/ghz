package jobs

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
			return
		}

		if err := UpdateJob(ctx, id, dsn, job); err != nil {
			c.JSON(http.StatusInternalServerError, UpdateResponse{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("update job(%v): %v", j.JobID, err),
			})
			return
		}
	}
}

func ListJobs(ctx context.Context, projectID, datasetName string) ([]dataset.WorkflowJob, error) {
	return make([]dataset.WorkflowJob, 0), nil
}

func UpdateJob(ctx context.Context, projectID, datasetName string, j *github.WorkflowJob) error {
	return nil
}
