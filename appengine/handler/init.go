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
	id, dsn := dataset.Name(owner, repository)

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
		view.FrequencyRunsMeta(id, dsn),
		view.FrequencyJobsMeta(id, dsn),
		view.PullReqsMeta(id, dsn),
		view.PullReqsLeadTimeMeta(id, dsn),
		view.PullReqsLeadTimeMedianMeta(id, dsn),
		view.PullReqsTTRMeta(id, dsn),
		view.PullReqsTTRMedianMeta(id, dsn),
		view.PullReqsFailureRate(id, dsn),
		view.PushedMeta(id, dsn),
		view.PushedLeadTimeMeta(id, dsn),
		view.PushedLeadTimeMedianMeta(id, dsn),
		view.PushedTTRMeta(id, dsn),
		view.PushedTTRMedianMeta(id, dsn),
		view.PushedFailureRate(id, dsn),
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
