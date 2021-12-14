package commits

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/pullreqs/commits"
)

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()
	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.PullReqCommitsTableMeta); err != nil {
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

	prs, err := GetPullReqs(ctx, datasetName, token)
	if err != nil {
		log.Printf("get pull requests: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	for _, p := range prs {
		list, err := commits.Fetch(ctx,
			&commits.FetchInput{
				Owner:      owner,
				Repository: repository,
				PAT:        os.Getenv("PAT"),
				Page:       0,
				PerPage:    100,
			},
			int(p.Number),
		)
		if err != nil {
			log.Printf("fetch: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		items := make([]interface{}, 0)
		for _, r := range list {
			message := regexpnl.ReplaceAllString(r.Commit.GetMessage(), " ")
			if len(message) > 64 {
				message = message[0:64]
			}

			items = append(items, dataset.PullReqCommits{
				Owner:      owner,
				Repository: repository,
				ID:         p.ID,
				Number:     p.Number,
				SHA:        r.GetSHA(),
				Login:      r.Commit.Author.GetName(),
				Date:       r.Commit.Author.GetDate(),
				Message:    message,
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

func GetPullReqs(ctx context.Context, datasetName string, nextToken int64) ([]PullReq, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where id > %v", table, nextToken)

	prs := make([]PullReq, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		prs = append(prs, PullReq{
			ID:     values[0].(int64),
			Number: values[1].(int64),
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return prs, nil
}

func NextToken(ctx context.Context, datasetName string) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqCommitsTableMeta.Name)
	query := fmt.Sprintf("select max(id), max(number) from `%v` limit 1", table)

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
