package runs

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
	"github.com/itsubaki/ghz/pkg/actions/runs"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	logf      = logger.Factory
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

	if _, err := runs.Fetch(ctx,
		&runs.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			LastID:     token,
		},
		func(list []*github.WorkflowRun) error {
			items := make([]interface{}, 0)
			for _, r := range list {
				items = append(items, dataset.WorkflowRun{
					Owner:        owner,
					Repository:   repository,
					WorkflowID:   r.GetWorkflowID(),
					WorkflowName: r.GetName(),
					RunID:        r.GetID(),
					RunNumber:    int64(r.GetRunNumber()),
					Status:       r.GetStatus(),
					Conclusion:   r.GetConclusion(),
					CreatedAt:    r.GetCreatedAt().Time,
					UpdatedAt:    r.GetUpdatedAt().Time,
					HeadSHA:      r.GetHeadSHA(),
				})
			}

			if err := dataset.Insert(ctx, dsn, dataset.WorkflowRunsMeta.Name, items); err != nil {
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

func GetNextToken(ctx context.Context, projectID, dsn string) (int64, int64, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, dsn, dataset.WorkflowRunsMeta.Name)
	query := fmt.Sprintf("select max(run_id), max(run_number) from `%v` limit 1", table)

	var rid, num int64
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 2 {
			return
		}

		if values[0] == nil || values[1] == nil {
			return
		}

		rid = values[0].(int64)
		num = values[1].(int64)
	}); err != nil {
		return -1, -1, fmt.Errorf("query(%v): %v", query, err)
	}

	return rid, num, nil
}
