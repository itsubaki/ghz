package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func LeadTimeWorkflowMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime_workflow",
		ViewQuery: fmt.Sprintf(
			`SELECT
				B.owner,
				B.repository,
				B.workflow_id,
                B.workflow_name,
				A.id as pullreq_id,
				A.number as pullreq_number,
				A.login,
			    A.title,
				A.message,
				A.sha as commit_sha,
				A.committed_at,
				B.updated_at as completed_at,
				TIMESTAMP_DIFF(B.updated_at, A.committed_at, MINUTE) as lead_time
			FROM %v as A
			INNER JOIN %v as B
			ON A.merge_commit_sha = B.head_sha
			WHERE B.conclusion = "success"`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, LeadTimePullReqsMeta(projectID, datasetName).Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.WorkflowRunsMeta.Name),
		),
	}
}
