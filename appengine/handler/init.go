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

	owner := c.Param("owner")
	repository := c.Param("repository")
	renew := c.Query("renew")
	traceID := c.GetString("trace_id")
	spanID := c.GetString("span_id")

	projectID := dataset.ProjectID
	dsn := dataset.Name(owner, repository)

	// ctx, err := tracer.NewContext(context.Background(), traceID, spanID)
	// if err != nil {
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	log := logger.New(projectID, traceID).NewReport(ctx)
	log.Debug("trace_id: %v, span_id: %v", traceID, spanID)

	// tra, err := tracer.New(projectID, "func Init")
	// if err != nil {
	// 	log.ErrorAndReport(c.Request, "new tracer: %v", err)
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	return
	// }

	// func(ctx context.Context) {
	// 	_, span := tra.Start(ctx, "Hello World")
	// 	defer span.End()

	// 	time.Sleep(3 * time.Second)
	// }(ctx)

	// func(ctx context.Context) {
	// 	_, span := tra.Start(ctx, "FOOBAR")
	// 	defer span.End()

	// 	time.Sleep(1 * time.Second)
	// }(ctx)

	if strings.ToLower(renew) == "true" {
		if err := dataset.DeleteAllView(ctx, dsn); err != nil {
			log.ErrorAndReport(c.Request, "delete all view: %v", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		log.Debug("delete all view")
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
	log.Debug("created table/view")

	c.JSON(http.StatusOK, gin.H{
		"path": c.Request.URL.Path,
	})
}
