package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func LeadTimePullReqsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime_via_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.workflow_id,
				A.workflow_name,
				B.id,
				B.number,
				B.login,
				B.title,
				B.message,
				B.merge_commit_sha,
				B.sha,
				B.committed_at,
				A.updated_at as completed_at,
				TIMESTAMP_DIFF(A.updated_at, B.committed_at, MINUTE) as lead_time
			FROM %v as A
			INNER JOIN %v as B
			ON A.head_sha = B.merge_commit_sha
			WHERE A.conclusion = "success"
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.WorkflowRunsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, PullReqsMeta(projectID, datasetName).Name),
		),
	}
}
