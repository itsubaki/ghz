package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Incident struct {
	Owner       string    `bigquery:"owner" json:"owner"`
	Repository  string    `bigquery:"repository" json:"repository"`
	Description string    `bigquery:"description" json:"description"`
	SHA         string    `bigquery:"sha" json:"sha"`
	ResolvedAt  time.Time `bigquery:"resolved_at" json:"resolved_at"`
}

var IncidentsMeta = bigquery.TableMetadata{
	Name: "incidents",
	TimePartitioning: &bigquery.TimePartitioning{
		Type:  bigquery.MonthPartitioningType,
		Field: "resolved_at",
	},
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "description", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType, Required: true},
		{Name: "resolved_at", Type: bigquery.TimestampFieldType, Required: true},
	},
}
