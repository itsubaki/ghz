package pullreqs

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/appengine/dataset/view"
	"github.com/itsubaki/ghstats/pkg/pullreqs"
)

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, dsn, []bigquery.TableMetadata{
		dataset.PullReqsMeta,
		view.PullReqsMeta(id, dsn),
		view.LeadTimePullReqsMeta(id, dsn),
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := NextToken(ctx, id, dsn)
	if err != nil {
		log.Printf("get lastID: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, token)

	if _, err := pullreqs.Fetch(ctx,
		&pullreqs.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			State:      "all",
			LastID:     token,
		},
		func(list []*github.PullRequest) error {
			items := make([]interface{}, 0)
			for _, r := range list {
				title := regexpnl.ReplaceAllString(r.GetTitle(), " ")
				if len(title) > 64 {
					title = title[0:64]
				}

				items = append(items, dataset.PullReq{
					Owner:          owner,
					Repository:     repository,
					ID:             r.GetID(),
					Number:         int64(r.GetNumber()),
					Login:          r.User.GetLogin(),
					Title:          title,
					State:          r.GetState(),
					CreatedAt:      r.GetCreatedAt(),
					UpdatedAt:      r.GetUpdatedAt(),
					MergedAt:       r.GetMergedAt(),
					ClosedAt:       r.GetClosedAt(),
					MergeCommitSHA: r.GetMergeCommitSHA(),
				})
			}

			if err := dataset.Insert(ctx, dsn, dataset.PullReqsMeta.Name, items); err != nil {
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

func NextToken(ctx context.Context, projectID, datasetName string) (int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.PullReqsMeta.Name)
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
