package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Incident struct {
	Owner         string    `bigquery:"owner" json:"owner"`
	Repository    string    `bigquery:"repository" json:"repository"`
	ID            string    `bigquery:"id" json:"id"`
	Description   string    `bigquery:"description" json:"description"`
	PullReqNumber int64     `bigquery:"pullreq_number" json:"pullreq_number"`
	CreatedAt     time.Time `bigquery:"created_at" json:"created_at"`
	ResolvedAt    time.Time `bigquery:"resolved_at" json:"resolved_at"`
}

var IncidentMeta = bigquery.TableMetadata{
	Name: "incidents",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.StringFieldType, Required: true},
		{Name: "description", Type: bigquery.StringFieldType, Required: true},
		{Name: "pullreq_number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "resolved_at", Type: bigquery.TimestampFieldType, Required: true},
	},
}
