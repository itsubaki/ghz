package dataset

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

var projectID = func() string {
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
	return projectID, fmt.Sprintf("%v_%v", own, rep)
}

type Client struct {
	client *bigquery.Client
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

	return &Client{
		client: client,
	}, nil
}

func (c *Client) CreateIfNotExists(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	location := "US"
	if len(os.Getenv("DATASET_LOCATION")) > 0 {
		location = os.Getenv("DATASET_LOCATION")
	}

	if _, err := c.client.Dataset(datasetName).Metadata(ctx); err != nil {
		// not found then create dataset
		if err := c.client.Dataset(datasetName).Create(ctx, &bigquery.DatasetMetadata{
			Location: location,
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

func (c *Client) Delete(ctx context.Context, projectID, datasetName string, tableName []string) error {
	for _, n := range tableName {
		q := fmt.Sprintf("DELETE From `%v.%v.%v` WHERE true", projectID, datasetName, n)
		if err := c.Query(ctx, q, func(values []bigquery.Value) {}); err != nil {
			return fmt.Errorf("query(%v): %v", q, err)
		}
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	if err := c.client.Dataset(datasetName).Table(tableName).Inserter().Put(ctx, items); err != nil {
		return fmt.Errorf("insert %v/%v: %v", datasetName, tableName, err)
	}

	return nil
}

func (c *Client) Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	it, err := c.client.Query(query).Read(ctx)
	if err != nil {
		return fmt.Errorf("query(%v): %v", query, err)
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

func Delete(ctx context.Context, projectID, datasetName string, tableName []string) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Delete(ctx, projectID, datasetName, tableName)
}

func CreateIfNotExists(ctx context.Context, datasetName string, meta []bigquery.TableMetadata) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.CreateIfNotExists(ctx, datasetName, meta)
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
