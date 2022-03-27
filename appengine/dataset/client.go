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

var invalid = regexp.MustCompile(`[!?"'#$%&@\+\-\*/=~^;:,.|()\[\]{}<>]`)

func Name(owner, repository string) string {
	own := invalid.ReplaceAllString(owner, "_")
	rep := invalid.ReplaceAllString(repository, "_")
	return fmt.Sprintf("%v_%v", own, rep)
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

func (c *Client) Create(ctx context.Context, dsn string, meta []bigquery.TableMetadata) error {
	if _, err := c.client.Dataset(dsn).Metadata(ctx); err != nil {
		// not found then create dataset
		if err := c.client.Dataset(dsn).Create(ctx, &bigquery.DatasetMetadata{
			Location: c.location,
		}); err != nil {
			return fmt.Errorf("create %v: %v", dsn, err)
		}
	}

	for _, m := range meta {
		ref := c.client.Dataset(dsn).Table(m.Name)
		if _, err := ref.Metadata(ctx); err == nil {
			// already exists
			continue
		}

		if err := ref.Create(ctx, &m); err != nil {
			return fmt.Errorf("create %v/%v: %v", dsn, m.Name, err)
		}
	}

	return nil
}

func (c *Client) Delete(ctx context.Context, dsn string, tables []bigquery.TableMetadata) error {
	if _, err := c.client.Dataset(dsn).Metadata(ctx); err != nil {
		return fmt.Errorf("dataset(%v): %v", dsn, err)
	}

	for _, t := range tables {
		ref := c.client.Dataset(dsn).Table(t.Name)
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

func (c *Client) DeleteAllView(ctx context.Context, dsn string) error {
	if _, err := c.client.Dataset(dsn).Metadata(ctx); err != nil {
		return fmt.Errorf("dataset(%v): %v", dsn, err)
	}

	it := c.client.Dataset(dsn).Tables(ctx)
	for {
		t, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("table: %v", err)
		}
		m, err := t.Metadata(ctx)
		if err != nil {
			return fmt.Errorf("table metadata: %v", err)
		}
		if m.Type != bigquery.ViewTable {
			continue
		}

		if err := c.Delete(ctx, dsn, []bigquery.TableMetadata{*m}); err != nil {
			return fmt.Errorf("delete view: %v", err)
		}
	}

	return nil
}

func (c *Client) Insert(ctx context.Context, dsn, table string, items []interface{}) error {
	if err := c.client.Dataset(dsn).Table(table).Inserter().Put(ctx, items); err != nil {
		return fmt.Errorf("insert %v.%v.%v: %v", c.client.Project(), dsn, table, err)
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

func Create(ctx context.Context, dsn string, meta []bigquery.TableMetadata) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Create(ctx, dsn, meta)
}

func Delete(ctx context.Context, dsn string, tables []bigquery.TableMetadata) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Delete(ctx, dsn, tables)
}

func DeleteAllView(ctx context.Context, dsn string) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.DeleteAllView(ctx, dsn)
}

func Insert(ctx context.Context, dsn, table string, items []interface{}) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Insert(ctx, dsn, table, items)
}

func Query(ctx context.Context, query string, fn func(values []bigquery.Value)) error {
	c, err := New(ctx)
	if err != nil {
		return fmt.Errorf("new client: %v", err)
	}
	defer c.Close()

	return c.Query(ctx, query, fn)
}
