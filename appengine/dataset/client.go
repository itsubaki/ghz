package dataset

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"golang.org/x/xerrors"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
)

var ProjectID = func() string {
	creds, err := google.FindDefaultCredentials(context.Background())
	if err != nil {
		panic(fmt.Sprintf("find default credentials: %v", err))
	}

	return creds.ProjectID
}()

var invalid = regexp.MustCompile(`[!?"'#$%&@\+\-\*/=~^;:,.|()\[\]{}<>]`)

func Name(owner, repository string) (string, string) {
	own := invalid.ReplaceAllString(owner, "_")
	rep := invalid.ReplaceAllString(repository, "_")
	return ProjectID, fmt.Sprintf("%v_%v", own, rep)
}

type Client struct {
	client   *bigquery.Client
	location string
}

func New(ctx context.Context) (*Client, error) {
	creds, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return nil, fmt.Errorf("find default credentials: %v", err)
	}

	client, err := bigquery.NewClient(ctx, creds.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("new bigquery client: %v", err)
	}

	location := "US"
	if len(os.Getenv("DATASET_LOCATION")) > 0 {
		location = os.Getenv("DATASET_LOCATION")
	}

	return &Client{
		client:   client,
		location: location,
	}, nil
}

func (c *Client) Create(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	if _, err := c.client.Dataset(datasetName).Metadata(ctx); err != nil {
		// not found then create dataset
		if err := c.client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{
			Location: c.location,
		}); err != nil {
			return fmt.Errorf("create %v: %v", datasetName, err)
		}
	}

	for _, m := range meta {
		ref := c.client.Dataset(datasetName).Table(m.Name)
		if _, err := ref.Metadata(ctx); err == nil {
			// already exists
			continue
		}

		if err := ref.Create(ctx, &m); err != nil {
			return fmt.Errorf("create %v/%v: %v", datasetName, m.Name, err)
		}
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, datasetName string, tables []bigquery.TableMetadata) error {
	if _, err := c.client.Dataset(datasetName).Metadata(ctx); err != nil {
		return fmt.Errorf("dataset(%v): %v", datasetName, err)
	}

	for _, t := range tables {
		ref := c.client.Dataset(datasetName).Table(t.Name)
		if _, err := ref.Metadata(ctx); err != nil {
			// https://pkg.go.dev/cloud.google.com/go/bigquery#hdr-Errors
			var e *googleapi.Error
			if ok := xerrors.As(err, &e); ok && e.Code == http.StatusNotFound {
				// already deleted
				return nil
			}

			return fmt.Errorf("table(%v): %v", t.Name, err)
		}

		if err := ref.Delete(ctx); err != nil {
			return fmt.Errorf("delete table=%v: %v", t.Name, err)
		}
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	if err := c.client.Dataset(datasetName).Table(tableName).Inserter().Put(ctx, items); err != nil {
		return fmt.Errorf("insert %v.%v.%v: %v", ProjectID, datasetName, tableName, err)
	}

	return nil
}

func (c *Client) Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	q := c.client.Query(query)
	q.Location = c.location

	it, err := q.Read(ctx)
	if err != nil {
		return fmt.Errorf("query: %v", err)
	}

	var values []bigquery.Value
	for {
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}

		if err != nil {
			return fmt.Errorf("iterator: %v", err)
		}

		fn(values)
	}

	return nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) Raw() *bigquery.Client {
	return c.client
}

func Delete(ctx context.Context, datasetName string, tables []bigquery.TableMetadata) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Delete(ctx, datasetName, tables)
}

func Create(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Create(ctx, datasetName, meta)
}

func Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Insert(ctx, datasetName, tableName, items)
}

func Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Query(ctx, query, fn)
}
