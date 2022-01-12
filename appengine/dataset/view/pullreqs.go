package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PullReqsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.id,
				A.number,
				A.login,
				A.title,
				B.message,
				A.merge_commit_sha,
				B.sha,
				B.date as committed_at,
				A.merged_at,
				TIMESTAMP_DIFF(A.merged_at, B.date, MINUTE) as duration
			FROM %v as A
			INNER JOIN %v as B
			ON A.id = B.id
			WHERE A.state = "closed" AND A.merged_at != "0001-01-01 00:00:00 UTC"
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
		),
	}
}
