package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type WorkflowJob struct {
	WorkflowID   int64     `bigquery:"workflow_id"`
	WorkflowName string    `bigquery:"workflow_name"`
	RunID        int64     `bigquery:"run_id"`
	RunNumber    int       `bigquery:"run_number"`
	JobID        int64     `bigquery:"job_id"`
	JobName      string    `bigquery:"job_name"`
	Status       string    `bigquery:"status"`
	Conclusion   string    `bigquery:"conclusion"`
	StartedAt    time.Time `bigquery:"started_at"`
	CompletedAt  time.Time `bigquery:"completed_at"`
}

var WorkflowJobsTableMeta = bigquery.TableMetadata{
	Name: "workflow_jobs",
	Schema: bigquery.Schema{
		{Name: "workflow_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "workflow_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "run_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "run_number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "job_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "job_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "status", Type: bigquery.StringFieldType, Required: true},
		{Name: "conclusion", Type: bigquery.StringFieldType},
		{Name: "started_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "completed_at", Type: bigquery.TimestampFieldType},
	},
}
