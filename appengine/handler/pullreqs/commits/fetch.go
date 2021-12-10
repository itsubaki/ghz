package commits

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/pullreqs/commits"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	datasetName := dataset.Name(c.Query("owner"), c.Query("repository"))

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.PullReqCommitsTableMeta); err != nil {
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

	prs, err := GetPullReqs(ctx, datasetName, id)
	if err != nil {
		log.Printf("get pull requests: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := commits.FetchInput{
		Owner:      c.Query("owner"),
		Repository: c.Query("repository"),
		PAT:        os.Getenv("PAT"),
		Page:       0,
		PerPage:    100,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repository, id, number)

	for _, p := range prs {
		list, err := commits.Fetch(ctx, &in, int(p.Number))
		if err != nil {
			log.Printf("fetch: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		items := make([]interface{}, 0)
		for _, r := range list {
			items = append(items, dataset.PullReqCommits{
				Owner:      c.Query("owner"),
				Repository: c.Query("repository"),
				ID:         p.ID,
				Number:     p.Number,
				SHA:        r.GetSHA(),
				Login:      r.Commit.Author.GetName(),
				Date:       r.Commit.Author.GetDate(),
				Message:    strings.ReplaceAll(r.Commit.GetMessage(), "\n", " "),
			})
		}

		if err := dataset.Insert(ctx, datasetName, dataset.PullReqCommitsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

type PullReq struct {
	ID     int64
	Number int64
}

func GetPullReqs(ctx context.Context, datasetName string, lastID int64) ([]PullReq, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where id > %v", table, lastID)

	prs := make([]PullReq, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		prs = append(prs, PullReq{
			ID:     values[0].(int64),
			Number: values[1].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return prs, nil
}

func GetLastID(ctx context.Context, datasetName string) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqCommitsTableMeta.Name)
	query := fmt.Sprintf("select max(id), max(number) from `%v` limit 1", table)

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
