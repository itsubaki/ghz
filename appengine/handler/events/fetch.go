package events

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/events"
)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, []bigquery.TableMetadata{
		dataset.EventsMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, nil)

	if _, err := events.Fetch(ctx,
		&events.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
		},
		func(list []*github.Event) error {
			items := make([]interface{}, 0)
			for _, e := range list {
				if e.GetType() != "PushEvent" {
					continue
				}

				items = append(items, dataset.Event{
					Owner:      owner,
					Repository: repository,
					ID:         e.GetID(),
					Login:      e.GetActor().GetLogin(),
					Type:       e.GetType(),
					CreatedAt:  e.GetCreatedAt(),
					RawPayload: string(e.GetRawPayload()),
				})
			}

			if err := dataset.Insert(ctx, datasetName, dataset.EventsMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		}); err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}
