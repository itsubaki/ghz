package incidents

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
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

	ctx := context.Background()
	_, dsn := dataset.Name(in.Owner, in.Repository)

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.IncidentsMeta,
	}); err != nil {
		c.Error(err).SetMeta(gin.H{
			"message": fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	resolvedAt, err := time.Parse("2006-01-02 15:04:05 UTC", in.ResolvedAt)
	if err != nil {
		c.Error(err).SetMeta(gin.H{
			"message": fmt.Sprintf("parse time: %v", err),
		})
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
		c.Error(err).SetMeta(gin.H{
			"message": fmt.Sprintf("insert items: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, in)
}
