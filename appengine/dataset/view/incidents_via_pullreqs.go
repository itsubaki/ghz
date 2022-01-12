package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func IncidentsPullReqsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_incidents_via_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.description,
				B.id,
				B.number,
				B.login,
				B.title,
				B.message,
				B.merge_commit_sha,
				A.sha,
				B.merged_at,
				A.resolved_at,
				TIMESTAMP_DIFF(A.resolved_at, B.merged_at, MINUTE) as TTR
			FROM %v as A
			INNER JOIN %v as B
			ON A.sha = B.sha
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, PullReqsMeta(projectID, datasetName).Name),
		),
	}
}
