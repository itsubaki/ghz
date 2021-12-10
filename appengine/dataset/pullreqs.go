package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type PullReqs struct {
	Owner          string    `bigquery:"owner"`
	Repository     string    `bigquery:"repository"`
	ID             int64     `bigquery:"id"`
	Number         int64     `bigquery:"number"`
	Login          string    `bigquery:"login"`
	Title          string    `bigquery:"title"`
	State          string    `bigquery:"state"`
	CreatedAt      time.Time `bigquery:"created_at"`
	UpdatedAt      time.Time `bigquery:"updated_at"`
	MergedAt       time.Time `bigquery:"merged_at"`
	ClosedAt       time.Time `bigquery:"closed_at"`
	MergeCommitSHA string    `bigquery:"merge_commit_sha"`
}

var PullReqsTableMeta = bigquery.TableMetadata{
	Name: "pullreqs",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "title", Type: bigquery.StringFieldType, Required: true},
		{Name: "state", Type: bigquery.StringFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "updated_at", Type: bigquery.TimestampFieldType},
		{Name: "merged_at", Type: bigquery.TimestampFieldType},
		{Name: "closed_at", Type: bigquery.TimestampFieldType},
		{Name: "merge_commit_sha", Type: bigquery.StringFieldType, Required: true},
	},
}
