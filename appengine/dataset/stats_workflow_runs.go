package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type WorkflowRunStats struct {
	Owner        string    `bigquery:"owner"`
	Repository   string    `bigquery:"repository"`
	WorkflowID   int64     `bigquery:"workflow_id"`
	WorkflowName string    `bigquery:"workflow_name"`
	Start        time.Time `bigquery:"start"`
	End          time.Time `bigquery:"end"`
	RunsPerDay   float64   `bigquery:"runs_per_day"`
	FailureRate  float64   `bigquery:"failure_rate"`
	DurationAvg  float64   `bigquery:"duration_avg"`
	DurationVar  float64   `bigquery:"duration_var"`
}

var WorkflowRunStatsTableMeta = bigquery.TableMetadata{
	Name: "stats_workflow_run",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "workflow_id", Type: bigquery.StringFieldType, Required: true},
		{Name: "workflow_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "start", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "end", Type: bigquery.TimestampFieldType},
		{Name: "runs_per_day", Type: bigquery.FloatFieldType, Required: true},
		{Name: "failure_rate", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_avg", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_var", Type: bigquery.FloatFieldType, Required: true},
	},
}
