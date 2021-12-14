package dataset

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

type WorkflowJobStats struct {
	Owner        string     `bigquery:"owner"`
	Repository   string     `bigquery:"repository"`
	WorkflowID   int64      `bigquery:"workflow_id"`
	WorkflowName string     `bigquery:"workflow_name"`
	JobName      string     `bigquery:"job_name"`
	Year         int64      `bigquery:"year"`
	Week         int64      `bigquery:"week"`
	Start        civil.Date `bigquery:"start"`
	End          civil.Date `bigquery:"end"`
	RunsPerDay   float64    `bigquery:"runs_per_day"`
	FailureRate  float64    `bigquery:"failure_rate"`
	DurationAvg  float64    `bigquery:"duration_avg"`
	DurationVar  float64    `bigquery:"duration_var"`
}

var WorkflowJobStatsTableMeta = bigquery.TableMetadata{
	Name: "workflow_jobs_stats",
	TimePartitioning: &bigquery.TimePartitioning{
		Type:  bigquery.MonthPartitioningType,
		Field: "start",
	},
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "workflow_id", Type: bigquery.StringFieldType, Required: true},
		{Name: "workflow_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "job_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "year", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "week", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "start", Type: bigquery.DateFieldType, Required: true},
		{Name: "end", Type: bigquery.DateFieldType},
		{Name: "runs_per_day", Type: bigquery.FloatFieldType, Required: true},
		{Name: "failure_rate", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_avg", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_var", Type: bigquery.FloatFieldType, Required: true},
	},
}
