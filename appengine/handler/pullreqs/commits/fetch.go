package commits

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/pkg/pullreqs/commits"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
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

	token, _, err := GetNextToken(ctx, projectID, dsn)
	if err != nil {
		log.ErrorReport("get next token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("next token=%v", token)

	prs, err := ListPullReqs(ctx, projectID, dsn, token)
	if err != nil {
		log.ErrorReport("list pullreqs: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Debug("len(pullreqs)=%v", len(prs))

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
			log.ErrorReport("fetch: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
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

		if err := dataset.Insert(ctx, dsn, dataset.PullReqCommitsMeta.Name, items); err != nil {
			log.ErrorReport("insert items: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}

type PullReq struct {
	ID     int64
	Number int64
}

func ListPullReqs(ctx context.Context, projectID, dsn string, nextToken int64) ([]PullReq, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqsMeta.Name)
	query := fmt.Sprintf("select id, number from `%v` where id > %v", table, nextToken)

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

func GetNextToken(ctx context.Context, projectID, dsn string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.PullReqCommitsMeta.Name)
	query := fmt.Sprintf("select max(id), max(number) from `%v` limit 1", table)

	var pid, num int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		pid = values[0].(int64)
		num = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return pid, num, nil
}
