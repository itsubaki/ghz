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
	id, dsn := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, dsn, []bigquery.TableMetadata{
		dataset.PullReqsMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	open, err := GetPullReqs(ctx, id, dsn, "open")
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

		if err := UpdatePullReq(ctx, id, dsn, pr); err != nil {
			log.Printf("update pullreq: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := UpdatePullReqCommits(ctx, id, dsn, pr); err != nil {
			log.Printf("update pullreq commits: %v", err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	log.Println("updated")
	c.Status(http.StatusOK)
}

func UpdatePullReqCommits(ctx context.Context, projectID, datasetName string, r *github.PullRequest) error {
	return nil
}

func UpdatePullReq(ctx context.Context, projectID, datasetName string, r *github.PullRequest) error {
	client := dataset.New(ctx)
	defer client.Close()

	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.PullReqsMeta.Name)

	var query string
	if r.ClosedAt != nil {
		query = fmt.Sprintf("update %v set state = \"%v\", updated_at = \"%v\", merged_at = \"%v\", closed_at = \"%v\", merge_commit_sha = \"%v\" where id = %v",
			table,
			r.GetState(),
			r.GetUpdatedAt().Format("2006-01-02 15:04:05 UTC"),
			r.GetMergedAt().Format("2006-01-02 15:04:05 UTC"),
			r.GetClosedAt().Format("2006-01-02 15:04:05 UTC"),
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

func GetPullReqs(ctx context.Context, projectID, datasetName, state string) ([]dataset.PullReq, error) {
	client := dataset.New(ctx)
	defer client.Close()

	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where state = \"%v\"", table, state)

	out := make([]dataset.PullReq, 0)
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		out = append(out, dataset.PullReq{
			ID:     values[0].(int64),
			Number: values[1].(int64),
			State:  state,
		})
	}); err != nil {
		return nil, fmt.Errorf("query(%v): %v", query, err)
	}

	return out, nil
}
