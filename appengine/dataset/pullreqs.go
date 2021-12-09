package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type PullReqs struct {
	ID               int       `bigquery:"id"`
	Number           int       `bigquery:"number"`
	Login            string    `bigquery:"login"`
	Title            string    `bigquery:"title"`
	State            string    `bigquery:"state"`
	CreatedAt        time.Time `bigquery:"created_at"`
	UpdatedAt        time.Time `bigquery:"updated_at"`
	MergedAt         time.Time `bigquery:"merged_at"`
	ClosedAt         time.Time `bigquery:"closed_at"`
	merge_commit_sha string    `bigquery:"merge_commit_sha"`
}

var PullReqsTableMeta = bigquery.TableMetadata{
	Name: "pullreqs",
	Schema: bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "title", Type: bigquery.StringFieldType, Required: true},
		{Name: "state", Type: bigquery.StringFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "updated_at", Type: bigquery.TimestampFieldType},
		{Name: "merged_at", Type: bigquery.TimestampFieldType},
		{Name: "closed_at", Type: bigquery.TimestampFieldType},
		{Name: "merged_commit_sha", Type: bigquery.StringFieldType, Required: true},
	},
}
