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

	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.WorkflowRunsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, num, err := NextToken(ctx, datasetName)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v(%v)", c.Request.URL.Path, owner, repository, token, num)

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
					Owner:         owner,
					Repository:    repository,
					WorkflowID:    r.GetWorkflowID(),
					WorkflowName:  r.GetName(),
					RunID:         r.GetID(),
					RunNumber:     int64(r.GetRunNumber()),
					Status:        r.GetStatus(),
					Conclusion:    r.GetConclusion(),
					CreatedAt:     r.CreatedAt.Time,
					UpdatedAt:     r.UpdatedAt.Time,
					HeadCommitSHA: r.HeadCommit.GetID(),
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

func NextToken(ctx context.Context, datasetName string) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.WorkflowRunsTableMeta.Name)
	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v` limit 1", table)

	var id, num int64
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
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
