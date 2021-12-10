package runs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/actions/runs"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	datasetName := dataset.Name(c.Query("owner"), c.Query("repository"))

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.WorkflowRunsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	id, number, err := GetLastID(ctx, datasetName)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := runs.FetchInput{
		Owner:      c.Query("owner"),
		Repository: c.Query("repository"),
		PAT:        os.Getenv("PAT"),
		Page:       0,
		PerPage:    100,
		LastID:     id,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repository, in.LastID, number)

	if _, err := runs.Fetch(ctx, &in, func(list []*github.WorkflowRun) error {
		items := make([]interface{}, 0)
		for _, r := range list {
			items = append(items, dataset.WorkflowRun{
				Owner:         c.Query("owner"),
				Repository:    c.Query("repository"),
				WorkflowID:    r.GetWorkflowID(),
				WorkflowName:  r.GetName(),
				RunID:         r.GetID(),
				RunNumber:     int64(r.GetRunNumber()),
				Status:        r.GetStatus(),
				Conclusion:    r.GetConclusion(),
				CreatedAt:     r.CreatedAt.Time,
				UpdatedAt:     r.UpdatedAt.Time,
				HeadCommitSHA: *r.HeadCommit.ID,
			})
		}

		if err := dataset.Insert(ctx, datasetName, dataset.WorkflowRunsTableMeta.Name, items); err != nil {
			return fmt.Errorf("insert items: %v", err)
		}

		return nil
	}); err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetLastID(ctx context.Context, datasetName string) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowRunsTableMeta.Name)
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
