package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func LeadTimePullReqsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime_pullreqs",
		ViewQuery: fmt.Sprintf(
			`SELECT
				B.owner,
				B.repository,
				A.id,
				A.number,
				A.login,
				B.title,
				A.message,
				B.merge_commit_sha,
				A.sha,
				A.date as committed_at,
				B.merged_at,
				TIMESTAMP_DIFF(B.merged_at, A.date, MINUTE) as lead_time
			FROM %v as A
			INNER JOIN %v as B
			ON A.id = B.id
			WHERE B.merged_at != "0001-01-01 00:00:00 UTC"`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
		),
	}
}
