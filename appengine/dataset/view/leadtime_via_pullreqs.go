package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func LeadTimeWorkflowsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime_via_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					B.title,
					B.merge_commit_sha,
					A.id,
					A.number,
					A.login,
					A.message,
					A.sha,
					A.date,
				FROM %v as A
				INNER JOIN %v as B
				ON A.id = B.id
				WHERE B.merged_at != "0001-01-01 00:00:00 UTC"
			)
			SELECT
				B.owner,
				B.repository,
				B.workflow_id,
				B.workflow_name,
				A.id as pullreq_id,
				A.number as pullreq_number,
				A.login,
				A.title,
				A.message,
				A.merge_commit_sha,
				A.sha,
				A.date as committed_at,
				B.updated_at as completed_at,
				TIMESTAMP_DIFF(B.updated_at, A.date, MINUTE) as lead_time
			FROM A
			INNER JOIN %v as B
			ON A.merge_commit_sha = B.head_sha
			WHERE B.conclusion = "success"
			ORDER BY completed_at DESC
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.WorkflowRunsMeta.Name),
		),
	}
}
