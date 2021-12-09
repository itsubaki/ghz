package dataset

import (
	"context"
	"fmt"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
)

type Client struct {
	client *bigquery.Client
}

func NewClient(ctx context.Context) (*Client, error) {
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

func (c *Client) CreateIfNotExists(ctx context.Context, datasetName string, meta bigquery.TableMetadata) error {
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
