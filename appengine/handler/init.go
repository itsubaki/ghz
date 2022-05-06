package handler

import (
	"context"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/dataset/view"
	"github.com/itsubaki/ghz/appengine/logger"
	"github.com/itsubaki/ghz/appengine/tracer"
	"go.opentelemetry.io/otel"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	logf      = logger.Factory
	tra       = otel.Tracer("handler/init")
)

func Init(c *gin.Context) {
	owner := c.Param("owner")
	repository := c.Param("repository")
	renew := c.Query("renew")
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")
	traceTrue := c.GetBool("trace_true")

	ctx := context.Background()
	dsn := dataset.Name(owner, repository)
	log := logf.New(traceID, c.Request)
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
