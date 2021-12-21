package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type PushEvent struct {
	Owner      string    `bigquery:"owner"`
	Repository string    `bigquery:"repository"`
	ID         string    `bigquery:"id"`
	Login      string    `bigquery:"login"`
	Type       string    `bigquery:"type"`
	CreatedAt  time.Time `bigquery:"created_at"`
	HeadSHA    string    `bigquery:"head_sha"`
	SHA        string    `bigquery:"sha"`
	Message    string    `bigquery:"message"`
}

var EventsPushMeta = bigquery.TableMetadata{
	Name: "events_push",
	TimePartitioning: &bigquery.TimePartitioning{
		Type:  bigquery.MonthPartitioningType,
		Field: "created_at",
	},
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "type", Type: bigquery.StringFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "head_sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "message", Type: bigquery.StringFieldType, Required: true},
	},
}
