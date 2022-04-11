package runs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/actions/runs"
)

func Update(c *gin.Context) {
	ctx := context.Background()
	projectID := dataset.ProjectID

	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID).NewReport(ctx)

	list, err := ListRuns(ctx, projectID, dsn)
	if err != nil {
		log.ErrorAndReport(c.Request, "list jobs: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("runs=%v", list)

	for _, r := range list {
		run, err := runs.Get(ctx, &runs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			RunID:      r.RunID,
		})
		if err != nil {
			log.ErrorAndReport(c.Request, "get runID=%v: %v", r.RunID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := UpdateRun(ctx, projectID, dsn, run); err != nil {
			msg := strings.ReplaceAll(err.Error(), projectID, "$PROJECT_ID")
			log.Info("update runID=%v: %v", r.RunID, msg)
			continue
		}
		log.Debug("updated. runID=%v", r.RunID)
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func ListRuns(ctx context.Context, projectID, dsn string) ([]dataset.WorkflowRun, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowRunsMeta.Name)
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

func UpdateRun(ctx context.Context, projectID, dsn string, r *github.WorkflowRun) error {
	if r.GetStatus() != "completed" {
		return nil
	}

	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowRunsMeta.Name)
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
