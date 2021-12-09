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

	prs, err := GetPullRequests(ctx, id)
	if err != nil {
		log.Printf("get pull requests: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := commits.FetchInput{
		Owner:   c.Query("owner"),
		Repo:    c.Query("repository"),
		PAT:     os.Getenv("PAT"),
		Page:    0,
		PerPage: 100,
	}

	log.Printf("target=%v/%v, last_id=%v(%v)", in.Owner, in.Repo, id, number)

	for _, p := range prs {
		list, err := commits.Fetch(ctx, &in, int(p.Number))
		if err != nil {
			log.Printf("fetch: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		for _, r := range list {
			log.Printf("%v(%v)", r.GetSHA(), r.Commit.Author.GetDate())
		}

		items := make([]interface{}, 0)
		for _, r := range list {
			items = append(items, dataset.PullReqCommits{
				ID:      p.ID,
				Number:  p.Number,
				SHA:     r.GetSHA(),
				Login:   r.Commit.Author.GetName(),
				Date:    r.Commit.Author.GetDate(),
				Message: strings.ReplaceAll(r.Commit.GetMessage(), "\n", " "),
			})
		}

		if err := client.Insert(ctx, "raw", dataset.PullReqCommitsTableMeta.Name, items); err != nil {
			log.Printf("insert items: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

type PullRequest struct {
	ID     int64
	Number int64
}

func GetPullRequests(ctx context.Context, lastID int64) ([]PullRequest, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select id, number from `%v.%v.%v` where id > %v", client.ProjectID, "raw", dataset.PullReqsTableMeta.Name, lastID)
	it, err := client.Raw().Query(query).Read(ctx)
	if err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	prs := make([]PullRequest, 0)
	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("iterator: %v", err)
		}

		prs = append(prs, PullRequest{
			ID:     values[0].(int64),
			Number: values[1].(int64),
		})
	}

	return prs, nil
}

func GetLastID(ctx context.Context) (int64, int64, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return -1, -1, fmt.Errorf("new bigquery client: %v", err)
	}

	query := fmt.Sprintf("select max(id), max(number) from `%v.%v.%v` limit 1", client.ProjectID, "raw", dataset.PullReqCommitsTableMeta.Name)
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
