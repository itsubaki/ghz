package releases

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	_, dsn := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, dsn, []bigquery.TableMetadata{
		dataset.ReleasesMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}
