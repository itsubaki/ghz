package events

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghstats/appengine/dataset"
	"github.com/itsubaki/ghstats/pkg/events"
)

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	datasetName := dataset.Name(owner, repository)

	if err := dataset.CreateIfNotExists(ctx, datasetName, []bigquery.TableMetadata{
		dataset.EventsPushMeta,
		dataset.EventsMeta,
	}); err != nil {
		log.Printf("create if not exists: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	token, err := NextToken(ctx, datasetName)
	if err != nil {
		log.Printf("get lastSHA: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Printf("path=%v, target=%v/%v, next=%v", c.Request.URL.Path, owner, repository, token)

	if _, err := events.Fetch(ctx,
		&events.FetchInput{
			Owner:      owner,
			Repository: repository,
			PAT:        os.Getenv("PAT"),
			Page:       0,
			PerPage:    100,
			LastID:     token,
		},
		func(list []*github.Event) error {
			items := make([]interface{}, 0)
			for _, e := range list {
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
		},
		func(list []*github.Event) error {
			items := make([]interface{}, 0)
			for _, e := range list {
				if e.GetType() != "PushEvent" {
					continue
				}

				p := e.Payload().(*github.PushEvent)
				for _, c := range p.Commits {
					message := regexpnl.ReplaceAllString(c.GetMessage(), " ")
					if len(message) > 64 {
						message = message[0:64]
					}

					items = append(items, dataset.PushEvent{
						Owner:      owner,
						Repository: repository,
						ID:         e.GetID(),
						Login:      e.GetActor().GetLogin(),
						Type:       e.GetType(),
						CreatedAt:  e.GetCreatedAt(),
						HeadSHA:    p.GetHead(),
						SHA:        c.GetSHA(),
						Message:    message,
					})
				}
			}

			if err := dataset.Insert(ctx, datasetName, dataset.EventsPushMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		},
	); err != nil {
		log.Printf("fetch: %v", err)
		c.Status(http.StatusInternalServerError)
		return
	}

	log.Println("fetched")
	c.Status(http.StatusOK)
}

func NextToken(ctx context.Context, datasetName string) (string, error) {
	client := dataset.New(ctx)
	defer client.Close()

	table := fmt.Sprintf("%v.%v.%v", client.ProjectID, datasetName, dataset.EventsMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v`", table)

	var id string
	if err := client.Query(ctx, query, func(values []bigquery.Value) {
		if len(values) != 1 {
			return
		}

		if values[0] == nil {
			return
		}

		id = values[0].(string)
	}); err != nil {
		return "", fmt.Errorf("query(%v): %v", query, err)
	}

	return id, nil
}
