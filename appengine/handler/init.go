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
	"go.opentelemetry.io/otel/trace"
)

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
	log.DebugWith(spanID, "trace_id=%v, span_id=%v, trace_true=%v", traceID, spanID, traceTrue)

	tra, err := tracer.New(projectID)
	if err != nil {
		log.ErrorAndReport("new tracer: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer tra.ForceFlush(ctx)

	parent, err := tracer.NewContext(ctx, traceID, spanID, traceTrue)
	if err != nil {
		log.ErrorAndReport("new context: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := tra.Span(parent, "delete all view", func(child context.Context, s trace.Span) error {
		if strings.ToLower(renew) != "true" {
			s.AddEvent("renew flag is not true. nothing to do.")
			return nil
		}

		return dataset.DeleteAllView(child, dsn)
	}); err != nil {
		log.ErrorAndReport("delete all view: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err := tra.Span(parent, "create table/view", func(child context.Context, s trace.Span) error {
		return dataset.Create(child, dsn, []bigquery.TableMetadata{
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
	}); err != nil {
		log.ErrorAndReport("create if not exists: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}
