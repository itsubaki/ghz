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
			WITH A AS (
				SELECT
					A.owner,
					A.repository,
					A.title,
					B.message,
					A.merge_commit_sha,
					B.sha,
					B.date as committed_at,
					A.merged_at
				FROM %v as A
				INNER JOIN %v as B
				ON A.id = B.id
				WHERE A.state = "closed" AND A.merged_at != "0001-01-01 00:00:00 UTC"
			)
			SELECT
				A.owner,
				A.repository,
				A.title,
				A.message,
				B.description,
				A.merge_commit_sha,
				A.sha,
				A.merged_at,
				B.resolved_at,
				TIMESTAMP_DIFF(B.resolved_at, A.merged_at, MINUTE) as TTR
			FROM %v as B
			INNER JOIN A
			ON A.sha = B.sha
			ORDER BY A.merged_at DESC
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
