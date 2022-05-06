package pullreqs

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/pullreqs"
)

var (
	projectID = dataset.ProjectID
	logf      = logger.Factory
	regexpnl  = regexp.MustCompile(`\r\n|\r|\n`)
)

func Fetch(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logf.New(traceID, c.Request)

	token, err := GetNextToken(ctx, projectID, dsn)
	if err != nil {
		log.ErrorReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("next token=%v", token)

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
		},
	); err != nil {
		log.ErrorReport("fetch: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

func GetNextToken(ctx context.Context, projectID, dsn string) (int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v` limit 1", table)

	var pid int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		pid = values[0].(int64)
	}); err != nil {
		return -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return pid, nil
}
