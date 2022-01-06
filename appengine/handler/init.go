package handler

import (
	"context"
	"fmt"
	"net/http"

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

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.CommitsMeta,
		dataset.PullReqsMeta,
		dataset.PullReqCommitsMeta,
		dataset.WorkflowRunsMeta,
		dataset.WorkflowJobsMeta,
		dataset.EventsMeta,
		dataset.EventsPushMeta,
		dataset.ReleasesMeta,
		dataset.IncidentsMeta,
		view.IncidentsPullReqsMeta(id, dsn),
		view.IncidentsPushedMeta(id, dsn),
		view.LeadTimePullReqsMeta(id, dsn),
		view.LeadTimePushedMeta(id, dsn),
		view.WorkflowRunsMeta(id, dsn),
		view.WorkflowJobsMeta(id, dsn),
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
