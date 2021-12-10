package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Commits struct {
	Owner      string    `bigquery:"owner"`
	Repository string    `bigquery:"repository"`
	SHA        string    `bigquery:"sha"`
	Login      string    `bigquery:"login"`
	Date       time.Time `bigquery:"date"`
	Message    string    `bigquery:"message"`
}

var CommitsTableMeta = bigquery.TableMetadata{
	Name: "commits",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "date", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "message", Type: bigquery.StringFieldType, Required: true},
	},
}