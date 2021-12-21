package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Event struct {
	Owner      string    `bigquery:"owner"`
	Repository string    `bigquery:"repository"`
	ID         string    `bigquery:"id"`
	Login      string    `bigquery:"login"`
	Type       string    `bigquery:"type"`
	CreatedAt  time.Time `bigquery:"created_at"`
	RawPayload string    `bigquery:"raw_payload"`
}

var EventsMeta = bigquery.TableMetadata{
	Name: "events",
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
		{Name: "raw_payload", Type: bigquery.StringFieldType, Required: true},
	},
}
