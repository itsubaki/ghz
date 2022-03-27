package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func FrequencyJobsMeta(projectID, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_frequency_jobs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				workflow_id,
				workflow_name,
				job_name,
				DATE(completed_at) as date,
				COUNT(job_name) as runs,
				AVG(TIMESTAMP_DIFF(completed_at, started_at, MINUTE)) as duration_avg
			FROM %v
			WHERE conclusion = "success"
			GROUP BY owner, repository, workflow_id, workflow_name, job_name, date
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, dsn, dataset.WorkflowJobsMeta.Name),
		),
	}
}
