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
	"github.com/itsubaki/ghstats/appengine/dataset/view"
	"github.com/itsubaki/ghstats/pkg/actions/jobs"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, dsn, []bigquery.TableMetadata{
		dataset.WorkflowJobsMeta,
		view.WorkflowJobsMeta(id, dsn),
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, num, err := NextToken(ctx, id, dsn)
	if err != nil {
		log.Printf("get lastRunID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v(%v)", c.Request.URL.Path, owner, repository, token, num)

	runs, err := GetRuns(ctx, id, dsn, token)
	if err != nil {
		log.Printf("get runs: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

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
			log.Printf("fetch: %v", err)
			c.Status(http.StatusInternalServerError)
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
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetRuns(ctx context.Context, projectID, datasetName string, nextToken int64) ([]dataset.WorkflowRun, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowRunsMeta.Name)
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

func NextToken(ctx context.Context, projectID, datasetName string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.WorkflowJobsMeta.Name)
	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v` limit 1", table)

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
