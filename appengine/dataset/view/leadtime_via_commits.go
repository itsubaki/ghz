package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func LeadTimeCommitsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime_via_commits",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					A.login,
					A.message,
					A.sha,
					A.date
				FROM %v as A
				INNER JOIN %v as B
				ON A.sha = B.sha
			)
			SELECT
				B.owner,
				B.repository,
				B.workflow_id,
				B.workflow_name,
				A.login,
				A.message,
				A.sha,
				A.date as commited_at,
				B.updated_at as completed_at,
				TIMESTAMP_DIFF(B.updated_at, A.date, MINUTE) as lead_time
			FROM A
			INNER JOIN %v as B
			ON A.sha = B.head_sha
			WHERE B.conclusion = "success"
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.EventsPushMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.WorkflowRunsMeta.Name),
		),
	}
}
