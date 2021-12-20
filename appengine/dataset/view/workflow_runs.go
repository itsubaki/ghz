package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func WorkflowRunsMeta(projectID, datasetName, tableName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_workflow_runs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				workflow_id,
				workflow_name,
				DATE_ADD(DATE(created_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(created_at), INTERVAL -0 DAY)) +1 DAY) as week,
				count(workflow_id) / 7 as runs_per_day,
				AVG(TIMESTAMP_DIFF(updated_at, created_at,MINUTE)) as duration_avg,
				STDDEV(TIMESTAMP_DIFF(updated_at, created_at,MINUTE)) as duration_stddev
			FROM %v
			WHERE conclusion = "success"
			GROUP BY owner, repository, workflow_id, workflow_name, week
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, tableName),
		),
	}
}
