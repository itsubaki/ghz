package pullreqs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
	"google.golang.org/api/iterator"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	id, err := GetLastID(ctx)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := pullreqs.ListInput{
		Owner:      c.Query("owner"),
		Repository: c.Query("repository"),
		PAT:        os.Getenv("PAT"),
		Page:       0,
		PerPage:    100,
		State:      "all",
		LastID:     id,
	}
	log.Printf("target=%v/%v, last_id=%v", in.Owner, in.Repository, in.LastID)

	list, err := pullreqs.Fetch(ctx, &in)
	if err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, r := range list {
		log.Printf("%v(%v)", *r.ID, *r.Number)
	}

	items := make([]interface{}, 0)
	for _, r := range list {
		items = append(items, dataset.PullReqs{
			Owner:          c.Query("owner"),
			Repository:     c.Query("repository"),
			ID:             r.GetID(),
			Number:         r.GetNumber(),
			Login:          r.User.GetLogin(),
			Title:          r.GetTitle(),
			State:          r.GetState(),
			CreatedAt:      r.GetCreatedAt(),
			UpdatedAt:      r.GetUpdatedAt(),
			MergedAt:       r.GetMergedAt(),
			ClosedAt:       r.GetClosedAt(),
			MergeCommitSHA: r.GetMergeCommitSHA(),
		})
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := client.Insert(ctx, "raw", dataset.PullReqsTableMeta.Name, items); err != nil {
		log.Printf("insert items: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetLastID(ctx context.Context) (int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select max(id) from `%v.%v.%v` limit 1", client.ProjectID, "raw", dataset.PullReqsTableMeta.Name)
	it, err := client.Raw().Query(query).Read(ctx)
	if err != nil {
		return -1, fmt.Errorf("query(%v): %v", query, err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return -1, fmt.Errorf("iterator: %v", err)
		}
	}

	var id int64
	if len(values) > 0 && values[0] != nil {
		id = values[0].(int64)
	}

	return id, nil
}
