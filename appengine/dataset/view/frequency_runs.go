package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func FrequencyRunsMeta(projectID, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_frequency_runs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				workflow_id,
				workflow_name,
				DATE(updated_at) as date,
				COUNT(workflow_name) as runs,
				AVG(TIMESTAMP_DIFF(updated_at, created_at, MINUTE)) as duration_avg
			FROM %v
			WHERE conclusion = "success"
			GROUP BY owner, repository, workflow_id, workflow_name, date
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, dsn, dataset.WorkflowRunsMeta.Name),
		),
	}
}
