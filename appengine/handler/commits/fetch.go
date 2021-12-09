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
	"github.com/itsubaki/ghstats/pkg/commits"
	"google.golang.org/api/iterator"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	sha, err := GetLastSHA(ctx)
	if err != nil {
		log.Printf("get lastSHA: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	in := commits.ListInput{
		Owner:   c.Query("owner"),
		Repo:    c.Query("repository"),
		PAT:     os.Getenv("PAT"),
		Page:    0,
		PerPage: 100,
		LastSHA: sha,
	}

	list, err := commits.Fetch(ctx, &in)
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
		items = append(items, dataset.Commits{
			SHA:     r.GetSHA(),
			Login:   r.Commit.Author.GetName(),
			Date:    r.Commit.Author.GetDate(),
			Message: strings.ReplaceAll(r.Commit.GetMessage(), "\n", " "),
		})
	}

	client, err := dataset.New(ctx)
	if err != nil {
		log.Printf("new bigquery client: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if err := client.Insert(ctx, "raw", dataset.CommitsTableMeta.Name, items); err != nil {
		log.Printf("insert items: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func GetLastSHA(ctx context.Context) (string, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return "", fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, "raw", dataset.CommitsTableMeta.Name)
	query := fmt.Sprintf("select sha from `%v` where date = (select max(date) from `%v` limit 1)", table, table)
	it, err := client.Raw().Query(query).Read(ctx)
	if err != nil {
		return "", fmt.Errorf("query(%v): %v", query, err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return "", fmt.Errorf("iterator: %v", err)
		}
	}

	var sha string
	if len(values) > 0 && values[0] != nil {
		sha = values[0].(string)
	}

	return sha, nil
}
