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

func Update(c *gin.Context) {
	ctx := context.Background()
	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, dataset.PullReqsTableMeta); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	open, err := GetPullReqs(ctx, datasetName, "open")
	if err != nil {
		log.Printf("get pullreq with: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, len(open)=%v", c.Request.URL.Path, owner, repository, len(open))

	for _, r := range open {
		pr, err := pullreqs.Get(ctx, &pullreqs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Number:     int(r.Number),
		})
		if err != nil {
			log.Printf("get pullreq: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := UpdatePullReq(ctx, datasetName, pr); err != nil {
			log.Printf("update pullreq: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := UpdatePullReqCommits(ctx, datasetName, pr); err != nil {
			log.Printf("update pullreq commits: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("updated")
	c.Status(http.StatusOK)
}

func UpdatePullReqCommits(ctx context.Context, datasetName string, r *github.PullRequest) error {
	return nil
}

func UpdatePullReq(ctx context.Context, datasetName string, r *github.PullRequest) error {
	client, err := dataset.New(ctx)
	if err != nil {
		return fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)

	var query string
	if r.MergedAt != nil {
		query = fmt.Sprintf("update %v set state = \"%v\", updated_at = \"%v\", merged_at = \"%v\", merge_commit_sha = \"%v\" where id = %v",
			table,
			r.GetState(),
			r.UpdatedAt.Format("2006-01-02 15:04:05 UTC"),
			r.MergedAt.Format("2006-01-02 15:04:05 UTC"),
			r.GetMergeCommitSHA(),
			r.GetID(),
		)

	}

	if r.ClosedAt != nil {
		query = fmt.Sprintf("update %v set state = \"%v\", updated_at = \"%v\", closed_at = \"%v\", merge_commit_sha = \"%v\" where id = %v",
			table,
			r.GetState(),
			r.UpdatedAt.Format("2006-01-02 15:04:05 UTC"),
			r.ClosedAt.Format("2006-01-02 15:04:05 UTC"),
			r.GetMergeCommitSHA(),
			r.GetID(),
		)
	}

	if query == "" {
		return nil
	}

	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		return
	}); err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	return nil
}

func GetPullReqs(ctx context.Context, datasetName, state string) ([]dataset.PullReqs, error) {
	client, err := dataset.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.PullReqsTableMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where state = \"%v\"", table, state)

	out := make([]dataset.PullReqs, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		out = append(out, dataset.PullReqs{
			ID:     values[0].(int64),
			Number: values[1].(int64),
			State:  state,
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}
