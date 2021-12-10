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
	"google.golang.org/api/iterator"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	lastRunID, number, err := GetLastRunID(ctx)
	if err != nil {
		log.Printf("get lastRunID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	runs, err := GetRuns(ctx, lastRunID)
	if err != nil {
		log.Printf("get runs: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := jobs.FetchInput{
		Owner:   c.Query("owner"),
		Repo:    c.Query("repository"),
		PAT:     os.Getenv("PAT"),
		Page:    0,
		PerPage: 100,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repo, lastRunID, number)

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

		log.Printf("%v(%v)", r.RunID, r.RunNumber)

		items := make([]interface{}, 0)
		for _, j := range jobs {
			items = append(items, dataset.WorkflowJob{
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

		if err := client.Insert(ctx, "raw", dataset.WorkflowJobsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetRuns(ctx context.Context, lastID int64) ([]dataset.WorkflowRun, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select workflow_id, workflow_name, run_id, run_number from `%v.%v.%v` where run_id > %v", client.ProjectID, "raw", dataset.WorkflowRunsTableMeta.Name, lastID)
	it, err := client.Raw().Query(query).Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	runs := make([]dataset.WorkflowRun, 0)
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("iterator: %v", err)
		}

		runs = append(runs, dataset.WorkflowRun{
			WorkflowID:   values[0].(int64),
			WorkflowName: values[1].(string),
			RunID:        values[2].(int64),
			RunNumber:    int(values[3].(int64)),
		})
	}

	return runs, nil
}

func GetLastRunID(ctx context.Context) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v.%v.%v` limit 1", client.ProjectID, "raw", dataset.WorkflowJobsTableMeta.Name)
	it, err := client.Raw().Query(query).Read(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return -1, -1, fmt.Errorf("iterator: %v", err)
		}
	}

	var id int64
	if len(values) > 0 && values[0] != nil {
		id = values[0].(int64)
	}

	var number int64
	if len(values) > 1 && values[1] != nil {
		number = values[1].(int64)
	}

	return id, number, nil
}
