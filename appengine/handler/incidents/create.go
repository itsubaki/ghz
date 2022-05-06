package incidents

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/appengine/logger"
)

var (
	projectID = dataset.ProjectID
	logf      = logger.Factory
)

type Incident struct {
	Owner       string `json:"owner"`
	Repository  string `json:"repository"`
	Description string `json:"description"`
	SHA         string `json:"sha"`
	ResolvedAt  string `json:"resolved_at"`
}

func Create(c *gin.Context) {
	var in Incident
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("bind json: %v", err),
		})
		return
	}

	in.Owner = c.Param("owner")
	in.Repository = c.Param("repository")
	traceID := c.GetString("trace_id")

	ctx := context.Background()
	dsn := dataset.Name(in.Owner, in.Repository)
	log := logf.New(traceID, c.Request)

	resolvedAt, err := time.Parse("2006-01-02 15:04:05 UTC", in.ResolvedAt)
	if err != nil {
		log.ErrorReport("parse time: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	items := make([]interface{}, 0)
	items = append(items, dataset.Incident{
		Owner:       in.Owner,
		Repository:  in.Repository,
		Description: in.Description,
		SHA:         in.SHA,
		ResolvedAt:  resolvedAt,
	})

	if err := dataset.Insert(ctx, dsn, dataset.IncidentsMeta.Name, items); err != nil {
		log.ErrorReport("insert items: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, in)
}
