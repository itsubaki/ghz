package dataset

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
)

type PullReqStats struct {
	Owner        string     `bigquery:"owner"`
	Repository   string     `bigquery:"repository"`
	Week         int64      `bigquery:"week"`
	Start        civil.Date `bigquery:"start"`
	End          civil.Date `bigquery:"end"`
	MergedPerDay float64    `bigquery:"merged_per_day"`
	DurationAvg  float64    `bigquery:"duration_avg"` // avg(merged_timestamp - commit_timestamp)
	DurationVar  float64    `bigquery:"duration_var"` // avg(merged_timestamp - commit_timestamp)
}

var PullReqStatsTableMeta = bigquery.TableMetadata{
	Name: "pullreqs_stats",
	TimePartitioning: &bigquery.TimePartitioning{
		Type:  bigquery.MonthPartitioningType,
		Field: "start",
	},
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "week", Type: bigquery.IntegerFieldType, Required: true},
		{Name: "start", Type: bigquery.DateFieldType, Required: true},
		{Name: "end", Type: bigquery.DateFieldType},
		{Name: "merged_per_day", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_avg", Type: bigquery.FloatFieldType, Required: true},
		{Name: "duration_var", Type: bigquery.FloatFieldType, Required: true},
	},
}
