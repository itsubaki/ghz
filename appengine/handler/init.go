package handler

import (
	"context"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/dataset/view"
	"github.com/itsubaki/ghz/appengine/logger"
)

func Init(c *gin.Context) {
	ctx := context.Background()
	projectID := dataset.ProjectID

	owner := c.Param("owner")
	repository := c.Param("repository")
	traceID := c.GetString("trace_id")

	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID)

	if strings.ToLower(c.Query("renew")) == "true" {
		if err := dataset.DeleteAllView(ctx, dsn); err != nil {
			log.Error("delete all view: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.CommitsMeta,
		dataset.PullReqsMeta,
		dataset.PullReqCommitsMeta,
		dataset.EventsMeta,
		dataset.EventsPushMeta,
		dataset.ReleasesMeta,
		dataset.WorkflowRunsMeta,
		dataset.WorkflowJobsMeta,
		dataset.IncidentsMeta,
		view.FrequencyRunsMeta(projectID, dsn),
		view.FrequencyJobsMeta(projectID, dsn),
		view.PullReqsMeta(projectID, dsn),
		view.PullReqsLeadTimeMeta(projectID, dsn),
		view.PullReqsLeadTimeMedianMeta(projectID, dsn),
		view.PullReqsTTRMeta(projectID, dsn),
		view.PullReqsTTRMedianMeta(projectID, dsn),
		view.PullReqsFailureRate(projectID, dsn),
		view.PushedMeta(projectID, dsn),
		view.PushedLeadTimeMeta(projectID, dsn),
		view.PushedLeadTimeMedianMeta(projectID, dsn),
		view.PushedTTRMeta(projectID, dsn),
		view.PushedTTRMedianMeta(projectID, dsn),
		view.PushedFailureRate(projectID, dsn),
	}); err != nil {
		log.Error("create if not exists: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}
