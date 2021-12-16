package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func PullReqsLeadTimeMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_leadtime",
		ViewQuery: fmt.Sprintf(
			`SELECT
				A.id,
				A.number,
				A.login,
				A.message,
				B.merge_commit_sha,
				A.sha,
				B.merged_at,
				A.date as committed_at,
				TIMESTAMP_DIFF(B.merged_at, A.date, MINUTE) as lead_time
			FROM %v as A
			LEFT OUTER JOIN %v as B
			ON A.id = B.id
			WHERE B.merged_at != "0001-01-01 00:00:00 UTC"`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
		),
	}
}
