package events

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v40/github"
	"github.com/itsubaki/ghz/appengine/dataset"
	"github.com/itsubaki/ghz/pkg/events"
)

type Response struct {
	Path    string `json:"path"`
	Message string `json:"message,omitempty"`
}

var regexpnl = regexp.MustCompile(`\r\n|\r|\n`)

func Fetch(c *gin.Context) {
	ctx := context.Background()

	owner := c.Param("owner")
	repository := c.Param("repository")
	id, dsn := dataset.Name(owner, repository)

	token, err := NextToken(ctx, id, dsn)
	if err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("next token: %v", err),
		})
		return
	}

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

			if err := dataset.Insert(ctx, dsn, dataset.EventsMeta.Name, items); err != nil {
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

			if err := dataset.Insert(ctx, dsn, dataset.EventsPushMeta.Name, items); err != nil {
				return fmt.Errorf("insert items: %v", err)
			}

			return nil
		},
	); err != nil {
		c.Error(err).SetMeta(Response{
			Path:    c.Request.URL.Path,
			Message: fmt.Sprintf("fetch: %v", err),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Path: c.Request.URL.Path,
	})
}

func NextToken(ctx context.Context, projectID, datasetName string) (string, error) {
	table := fmt.Sprintf("%v.%v.%v", projectID, datasetName, dataset.EventsMeta.Name)
	query := fmt.Sprintf("select max(id) from `%v`", table)

	var id string
	if err := dataset.Query(ctx, query, func(values []bigquery.Value) {
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
