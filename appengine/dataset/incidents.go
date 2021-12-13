package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type Incident struct {
	Owner        string    `bigquery:"owner"`
	Repository   string    `bigquery:"repository"`
	ID           int64     `bigquery:"id"`
	DeploymentID int64     `bigquery:"deployment_id"`
	CreatedAt    time.Time `bigquery:"created_at"`
	ResolvedAt   time.Time `bigquery:"resolved_at"`
}

var IncidentTableMeta = bigquery.TableMetadata{
	Name: "incidents",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "deployment_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "resolved_at", Type: bigquery.TimestampFieldType, Required: true},
	},
}
