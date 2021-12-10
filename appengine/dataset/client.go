package dataset

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

type Client struct {
	client    *bigquery.Client
	ProjectID string
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
		client:    client,
		ProjectID: creds.ProjectID,
	}, nil
}

func (c *Client) CreateIfNotExists(ctx context.Context, datasetName string, meta bigquery.TableMetadata) error {
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

	ref := c.client.Dataset(datasetName).Table(meta.Name)
	if _, err := ref.Metadata(ctx); err == nil {
		// already exists
		return nil
	}

	if err := ref.Create(ctx, &meta); err != nil {
		return fmt.Errorf("create %v/%v: %v", datasetName, meta.Name, err)
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	if err := c.client.Dataset(datasetName).Table(tableName).Inserter().Put(ctx, items); err != nil {
		return fmt.Errorf("insert %v/%v: %v", datasetName, tableName, err)
	}

	return nil
}

func (c *Client) Raw() *bigquery.Client {
	return c.client
}

func CreateIfNotExists(ctx context.Context, datasetName string, meta bigquery.TableMetadata) error {
	client, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new bigquery client: %v", err)
	}

	return client.CreateIfNotExists(ctx, datasetName, meta)
}

func Insert(ctx context.Context, datasetName, tableName string, items []interface{}) error {
	client, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new bigquery client: %v", err)
	}

	return client.Insert(ctx, datasetName, tableName, items)
}

func Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	client, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new bigquery client: %v", err)
	}

	it, err := client.Raw().Query(query).Read(ctx)
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

func Name(owner, repository string) string {
	return fmt.Sprintf("%v_%v", owner, repository)
}
