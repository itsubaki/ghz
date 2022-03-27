package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/dataset/view"
)

type Response struct {
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

func Init(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	projectID := c.GetString("project_id")
	dsn := dataset.Name(owner, repository)

	if strings.ToLower(c.Query("renew")) == "true" {
		if err := dataset.DeleteAllView(ctx, dsn); err != nil {
			c.Error(err).SetMeta(Response{
				Path:    c.Request.URL.Path,
				Message: fmt.Sprintf("delete all view: %v", err),
			})
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
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Path: c.Request.URL.Path,
	})
}
