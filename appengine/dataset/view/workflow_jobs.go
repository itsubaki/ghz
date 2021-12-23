package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func WorkflowJobsMeta(projectID, datasetName, tableName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_workflow_jobs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				workflow_id,
				workflow_name,
				job_name,
				DATE_ADD(DATE(started_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(started_at), INTERVAL -0 DAY)) +1 DAY) as week,
				count(job_name) as runs,
				AVG(TIMESTAMP_DIFF(completed_at, started_at,MINUTE)) as duration_avg,
				STDDEV(TIMESTAMP_DIFF(completed_at, started_at,MINUTE)) as duration_stddev
			FROM %v
			WHERE conclusion = "success"
			GROUP BY owner, repository, workflow_id, workflow_name, job_name, week
			ORDER BY week DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, tableName),
		),
	}
}
