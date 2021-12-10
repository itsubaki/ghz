package jobs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/actions/jobs"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	datasetName := dataset.Name(c.Query("owner"), c.Query("repository"))

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.WorkflowJobsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	lastRunID, number, err := GetLastRunID(ctx, datasetName)
	if err != nil {
		log.Printf("get lastRunID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	runs, err := GetRuns(ctx, datasetName, lastRunID)
	if err != nil {
		log.Printf("get runs: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := jobs.FetchInput{
		Owner:      c.Query("owner"),
		Repository: c.Query("repository"),
		PAT:        os.Getenv("PAT"),
		Page:       0,
		PerPage:    100,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repository, lastRunID, number)

	for _, r := range runs {
		if r.RunID <= lastRunID {
			continue
		}

		jobs, err := jobs.Fetch(ctx, &in, r.RunID)
		if err != nil {
			log.Printf("fetch: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		items := make([]interface{}, 0)
		for _, j := range jobs {
			items = append(items, dataset.WorkflowJob{
				Owner:        in.Owner,
				Repository:   in.Repository,
				WorkflowID:   r.WorkflowID,
				WorkflowName: r.WorkflowName,
				RunID:        r.RunID,
				RunNumber:    r.RunNumber,
				JobID:        *j.ID,
				JobName:      *j.Name,
				Status:       *j.Status,
				Conclusion:   *j.Conclusion,
				StartedAt:    j.StartedAt.Time,
				CompletedAt:  j.CompletedAt.Time,
			})
		}

		if err := dataset.Insert(ctx, datasetName, dataset.WorkflowJobsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetRuns(ctx context.Context, datasetName string, lastID int64) ([]dataset.WorkflowRun, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowRunsTableMeta.Name)
	query := fmt.Sprintf("select workflow_id, workflow_name, run_id, run_number from `%v` where run_id > %v", table, lastID)

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

func GetLastRunID(ctx context.Context, datasetName string) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowJobsTableMeta.Name)
	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v` limit 1", table)

	var id, number int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		id = values[0].(int64)
		number = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return id, number, nil
}
