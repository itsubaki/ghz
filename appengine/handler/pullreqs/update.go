package pullreqs

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/pullreqs"
)

func Update(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID).NewReport(ctx, c.Request)

	open, err := ListPullReqs(ctx, projectID, dsn, "open")
	if err != nil {
		log.ErrorReport("list pullreqs: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("len(pullreqs)=%v", len(open))

	for _, r := range open {
		pr, err := pullreqs.Get(ctx, &pullreqs.GetInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Number:     int(r.Number),
		})
		if err != nil {
			log.ErrorReport("fetch pullreq: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if err := UpdatePullReq(ctx, projectID, dsn, pr); err != nil {
			log.Info("update pullreq(%v): %v", r.Number, err)
			continue
		}
		log.Debug("updated. pr=%v", r.Number)

		if err := UpdatePullReqCommits(ctx, projectID, dsn, pr); err != nil {
			log.Info("update commits(pr=%v): %v", r.Number, err)
			continue
		}

		log.Debug("updated commits. (pr=%v)", r.Number)
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func ListPullReqs(ctx context.Context, projectID, dsn, state string) ([]dataset.PullReq, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where state = \"%v\"", table, state)

	out := make([]dataset.PullReq, 0)
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
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

func UpdatePullReq(ctx context.Context, projectID, dsn string, r *github.PullRequest) error {
	if r.ClosedAt == nil {
		return nil
	}

	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("update %v set state = \"%v\", updated_at = \"%v\", merged_at = \"%v\", closed_at = \"%v\", merge_commit_sha = \"%v\" where id = %v",
		table,
		r.GetState(),
		r.GetUpdatedAt().Format("2006-01-02 15:04:05 UTC"),
		r.GetMergedAt().Format("2006-01-02 15:04:05 UTC"),
		r.GetClosedAt().Format("2006-01-02 15:04:05 UTC"),
		r.GetMergeCommitSHA(),
		r.GetID(),
	)

	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		return
	}); err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
	}

	return nil
}

func UpdatePullReqCommits(ctx context.Context, projectID, dsn string, r *github.PullRequest) error {
	return nil
}
