package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type PullReqCommits struct {
	ID      int64     `bigquery:"id"`
	Number  int64     `bigquery:"number"`
	SHA     string    `bigquery:"sha"`
	Login   string    `bigquery:"login"`
	Date    time.Time `bigquery:"date"`
	Message string    `bigquery:"message"`
}

var PullReqCommitsTableMeta = bigquery.TableMetadata{
	Name: "pullreqs_commits",
	Schema: bigquery.Schema{
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "login", Type: bigquery.StringFieldType, Required: true},
		{Name: "date", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "message", Type: bigquery.StringFieldType, Required: true},
	},
}
