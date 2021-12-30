package incidents

import (
	"context"
	"fmt"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func Create(c *gin.Context) {
	var in dataset.Incident
	if err := c.BindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("bind json: %v", err),
		})
		return
	}
	in.Owner = c.Param("owner")
	in.Repository = c.Param("repository")

	if in.ResolvedAt.Year() == 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("resolved_at(%v) is invalid", in.ResolvedAt),
		})
		return
	}

	ctx := context.Background()
	_, dsn := dataset.Name(in.Owner, in.Repository)

	if err := dataset.Create(ctx, dsn, []bigquery.TableMetadata{
		dataset.IncidentsMeta,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("create if not exists: %v", err),
		})
		return
	}

	items := make([]interface{}, 0)
	items = append(items, in)

	if err := dataset.Insert(ctx, dsn, dataset.IncidentsMeta.Name, items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("insert items: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, in)
}
