package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func WorkflowJobsMeta(projectID, datasetName string) bigquery.TableMetadata {
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
				AVG(TIMESTAMP_DIFF(completed_at, started_at,MINUTE)) as duration_avg
			FROM %v
			WHERE conclusion = "success"
			GROUP BY owner, repository, workflow_id, workflow_name, job_name, week
			ORDER BY week DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.WorkflowJobsMeta.Name),
		),
	}
}
