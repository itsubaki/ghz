package runs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/actions/runs"
	"google.golang.org/api/iterator"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	id, number, err := GetLastID(ctx)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := runs.FetchInput{
		Owner:   c.Query("owner"),
		Repo:    c.Query("repository"),
		PAT:     os.Getenv("PAT"),
		Page:    0,
		PerPage: 100,
		LastID:  id,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repo, in.LastID, number)

	list, err := runs.Fetch(ctx, &in)
	if err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, r := range list {
		log.Printf("%v(%v)", r.GetID(), r.GetRunNumber())
	}

	items := make([]interface{}, 0)
	for _, r := range list {
		items = append(items, dataset.WorkflowRun{
			WorkflowID:    r.GetWorkflowID(),
			WorkflowName:  r.GetName(),
			RunID:         r.GetID(),
			RunNumber:     r.GetRunNumber(),
			Status:        r.GetStatus(),
			Conclusion:    r.GetConclusion(),
			CreatedAt:     r.CreatedAt.Time,
			UpdatedAt:     r.UpdatedAt.Time,
			HeadCommitSHA: *r.HeadCommit.ID,
		})
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := client.Insert(ctx, "raw", dataset.WorkflowRunsTableMeta.Name, items); err != nil {
		log.Printf("insert items: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetLastID(ctx context.Context) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v.%v.%v` limit 1", client.ProjectID, "raw", dataset.WorkflowRunsTableMeta.Name)
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
