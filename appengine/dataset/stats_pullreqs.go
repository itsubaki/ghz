package dataset

import (
	"time"

	"cloud.google.com/go/bigquery"
)

type PullReqStats struct {
	Owner         string    `bigquery:"owner"`
	Repository    string    `bigquery:"repository"`
	Start         time.Time `bigquery:"start"`
	End           time.Time `bigquery:"end"`
	CreatedPerDay float64   `bigquery:"created_per_day"`
	MergedPerDay  float64   `bigquery:"merged_per_day"`
	DurationAvg   float64   `bigquery:"duration_avg"` // merged_timestamp - first_commit_timestamp
	DurationVar   float64   `bigquery:"duration_var"`
}

var PullReqStatsTableMeta = bigquery.TableMetadata{
	Name: "stats_pullreqs",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "start", Type: bigquery.TimestampFieldType, Required: true},
		{Name: "end", Type: bigquery.TimestampFieldType},
		{Name: "created_per_day", Type: bigquery.FloatFieldType, Required: true},
		{Name: "merged_per_day", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_avg", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_var", Type: bigquery.FloatFieldType, Required: true},
	},
}
