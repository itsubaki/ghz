package pullreqs

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
	"github.com/itsubaki/ghstats/pkg/pullreqs"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	datasetName := dataset.Name(c.Query("owner"), c.Query("repository"))

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.PullReqsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	id, err := GetLastID(ctx, datasetName)
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

	if _, err := pullreqs.Fetch(ctx, &in, func(list []*github.PullRequest) error {
		items := make([]interface{}, 0)
		for _, r := range list {
			items = append(items, dataset.PullReqs{
				Owner:          c.Query("owner"),
				Repository:     c.Query("repository"),
				ID:             r.GetID(),
				Number:         int64(r.GetNumber()),
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

		if err := dataset.Insert(ctx, datasetName, dataset.PullReqsTableMeta.Name, items); err != nil {
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

func GetLastID(ctx context.Context, datasetName string) (int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v` limit 1", table)

	var id int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		id = values[0].(int64)
	}); err != nil {
		return -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return id, nil
}
