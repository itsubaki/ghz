package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type DORAStats struct {
	Owner                string    `bigquery:"owner"`
	Repository           string    `bigquery:"repository"`
	Start                time.Time `bigquery:"start"`
	End                  time.Time `bigquery:"end"`
	DeploymentFrequency  float64   `bigquery:"deployment_frequency"`    // deployment per day in production
	ChangeFailureRate    float64   `bigquery:"change_failure_rate"`     // failure count per deployment
	LeadTimeForChanges   time.Time `bigquery:"lead_time_for_changes"`   // deployed_timestamp - commit_timestamp
	TimeToRestoreService time.Time `bigquery:"time_to_restore_service"` // resolved_timestamp - created_timestamp
}

var DORAStatsTableMeta = bigquery.TableMetadata{
	Name: "stats_dora",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "start", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "end", Type: bigquery.TimestampFieldType},
		{Name: "deployment_frequency", Type: bigquery.FloatFieldType, Required: true},
		{Name: "change_failure_rate", Type: bigquery.FloatFieldType},
		{Name: "lead_time_for_changes", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "time_to_restore_service", Type: bigquery.TimestampFieldType},
	},
}
