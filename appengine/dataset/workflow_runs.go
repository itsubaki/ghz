package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type WorkflowRun struct {
	Owner         string    `bigquery:"onwer"`
	Repository    string    `bigquery:"repository"`
	WorkflowID    int64     `bigquery:"workflow_id"`
	WorkflowName  string    `bigquery:"workflow_name"`
	RunID         int64     `bigquery:"run_id"`
	RunNumber     int       `bigquery:"run_number"`
	Status        string    `bigquery:"status"`
	Conclusion    string    `bigquery:"conclusion"`
	CreatedAt     time.Time `bigquery:"created_at"`
	UpdatedAt     time.Time `bigquery:"updated_at"`
	HeadCommitSHA string    `bigquery:"head_commit_sha"`
}

var WorkflowRunsTableMeta = bigquery.TableMetadata{
	Name: "workflow_runs",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "workflow_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "workflow_name", Type: bigquery.StringFieldType, Required: true},
		{Name: "run_id", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "run_number", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "status", Type: bigquery.StringFieldType, Required: true},
		{Name: "conclusion", Type: bigquery.StringFieldType},
		{Name: "created_at", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "updated_at", Type: bigquery.TimestampFieldType},
		{Name: "head_commit_sha", Type: bigquery.StringFieldType, Required: true},
	},
}
