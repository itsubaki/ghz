package jobs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/actions/jobs"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	projectID := dataset.ProjectID

	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID).NewReport(ctx, c.Request)

	token, _, err := GetNextToken(ctx, projectID, dsn)
	if err != nil {
		log.ErrorAndReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("next token=%v", token)

	runs, err := ListRuns(ctx, projectID, dsn, token)
	if err != nil {
		log.ErrorAndReport("list runs: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("len(runs)=%v", len(runs))

	for _, r := range runs {
		jobs, err := jobs.Fetch(ctx,
			&jobs.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
			},
			r.RunID,
		)
		if err != nil {
			log.ErrorAndReport("fetch runID=%v: %v", r.RunID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		items := make([]interface{}, 0)
		for _, j := range jobs {
			items = append(items, dataset.WorkflowJob{
				Owner:        owner,
				Repository:   repository,
				WorkflowID:   r.WorkflowID,
				WorkflowName: r.WorkflowName,
				RunID:        r.RunID,
				RunNumber:    r.RunNumber,
				JobID:        j.GetID(),
				JobName:      j.GetName(),
				Status:       j.GetStatus(),
				Conclusion:   j.GetConclusion(),
				StartedAt:    j.GetStartedAt().Time,
				CompletedAt:  j.GetCompletedAt().Time,
			})
		}

		if err := dataset.Insert(ctx, dsn, dataset.WorkflowJobsMeta.Name, items); err != nil {
			log.ErrorAndReport("insert items runID=%v: %v", r.RunID, err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func ListRuns(ctx context.Context, projectID, dsn string, nextToken int64) ([]dataset.WorkflowRun, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowRunsMeta.Name)
	query := fmt.Sprintf("select workflow_id, workflow_name, run_id, run_number from `%v` where run_id > %v", table, nextToken)

	runs := make([]dataset.WorkflowRun, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		runs = append(runs, dataset.WorkflowRun{
			WorkflowID:   values[0].(int64),
			WorkflowName: values[1].(string),
			RunID:        values[2].(int64),
			RunNumber:    values[3].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return runs, nil
}

func GetNextToken(ctx context.Context, projectID, dsn string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowJobsMeta.Name)
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
