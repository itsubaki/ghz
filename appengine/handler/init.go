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
	renew := c.Query("renew")
	traceID := c.GetString("trace_id")

	dsn := dataset.Name(owner, repository)
	log := logger.New(projectID, traceID).NewReport(ctx)

	// exporter, err := trace.New(trace.WithProjectID(projectID))
	// if err != nil {
	// 	log.ErrorAndReport(c.Request, "texporter.NewExporter: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }
	// provider := sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter))
	// defer provider.ForceFlush(ctx)

	// otel.SetTracerProvider(provider)
	// tracer := otel.GetTracerProvider().Tracer(c.Request.URL.Path)

	// if err := func(ctx context.Context) error {
	// 	ctx, span := tracer.Start(ctx, "delete all view")
	// 	defer span.End()

	// 	if strings.ToLower(renew) != "true" {
	// 		return nil
	// 	}

	// 	return dataset.DeleteAllView(ctx, dsn)
	// }(ctx); err != nil {
	// 	log.ErrorAndReport(c.Request, "delete all view: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// if err := func(ctx context.Context) error {
	// 	ctx, span := tracer.Start(ctx, "create table/view")
	// 	defer span.End()

	// 	return dataset.Create(ctx, dsn, []bigquery.TableMetadata{
	// 		dataset.CommitsMeta,
	// 		dataset.PullReqsMeta,
	// 		dataset.PullReqCommitsMeta,
	// 		dataset.EventsMeta,
	// 		dataset.EventsPushMeta,
	// 		dataset.ReleasesMeta,
	// 		dataset.WorkflowRunsMeta,
	// 		dataset.WorkflowJobsMeta,
	// 		dataset.IncidentsMeta,
	// 		view.FrequencyRunsMeta(dsn),
	// 		view.FrequencyJobsMeta(dsn),
	// 		view.PullReqsMeta(dsn),
	// 		view.PullReqsLeadTimeMeta(dsn),
	// 		view.PullReqsLeadTimeMedianMeta(dsn),
	// 		view.PullReqsTTRMeta(dsn),
	// 		view.PullReqsTTRMedianMeta(dsn),
	// 		view.PullReqsFailureRate(dsn),
	// 		view.PushedMeta(dsn),
	// 		view.PushedLeadTimeMeta(dsn),
	// 		view.PushedLeadTimeMedianMeta(dsn),
	// 		view.PushedTTRMeta(dsn),
	// 		view.PushedTTRMedianMeta(dsn),
	// 		view.PushedFailureRate(dsn),
	// 	})
	// }(ctx); err != nil {
	// 	log.ErrorAndReport(c.Request, "create if not exists: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	if strings.ToLower(renew) == "true" {
		if err := dataset.DeleteAllView(ctx, dsn); err != nil {
			log.ErrorAndReport(c.Request, "delete all view: %v", err)
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
	}); err != nil {
		log.ErrorAndReport(c.Request, "create if not exists: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}
