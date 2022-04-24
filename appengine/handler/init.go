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
	"github.com/itsubaki/ghz/appengine/tracer"
	"go.opentelemetry.io/otel"
)

var tra = otel.Tracer("handler/init")

func Init(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	renew := c.Query("renew")
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")
	traceTrue := c.GetBool("trace_true")

	projectID := dataset.ProjectID
	dsn := dataset.Name(owner, repository)

	log := logger.New(projectID, traceID).NewReport(ctx, c.Request)
	log.SpanOf(spanID).Debug("trace=%v", traceTrue)

	parent, err := tracer.NewContext(ctx, traceID, spanID, traceTrue)
	if err != nil {
		log.ErrorReport("new context: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := func() error {
		c, s := tra.Start(parent, "delete all view")
		defer s.End()

		if strings.ToLower(renew) != "true" {
			s.AddEvent("renew flag is not true. nothing to do.")
			return nil
		}

		return dataset.DeleteAllView(c, dsn)
	}(); err != nil {
		log.ErrorReport("delete all view: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := func() error {
		c, s := tra.Start(parent, "create table/view")
		defer s.End()

		return dataset.Create(c, dsn, []bigquery.TableMetadata{
			dataset.CommitsMeta,
			dataset.PullReqsMeta,
			dataset.PullReqCommitsMeta,
			dataset.EventsMeta,
			dataset.EventsPushMeta,
			dataset.ReleasesMeta,
			dataset.WorkflowRunsMeta,
			dataset.WorkflowJobsMeta,
			dataset.IncidentsMeta,
			view.FrequencyRunsMeta(dsn),
			view.FrequencyJobsMeta(dsn),
			view.PullReqsMeta(dsn),
			view.PullReqsLeadTimeMeta(dsn),
			view.PullReqsLeadTimeMedianMeta(dsn),
			view.PullReqsTTRMeta(dsn),
			view.PullReqsTTRMedianMeta(dsn),
			view.PullReqsFailureRate(dsn),
			view.PushedMeta(dsn),
			view.PushedLeadTimeMeta(dsn),
			view.PushedLeadTimeMedianMeta(dsn),
			view.PushedTTRMeta(dsn),
			view.PushedTTRMedianMeta(dsn),
			view.PushedFailureRate(dsn),
		})
	}(); err != nil {
		log.ErrorReport("create if not exists: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}
