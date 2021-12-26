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
					owner,
					repository,
					count(merged_at) as merged,
					DATE(merged_at) as date
				FROM %v
				WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
				GROUP BY owner, repository, date
			), B AS (
				SELECT
					Date(A.date) as date,
					count(A.date) as failure,
					avg(TIMESTAMP_DIFF(B.resolved_at, A.date, MINUTE)) as MTTR					
				FROM %v as A
				INNER JOIN %v as B
				ON A.sha = B.sha
				GROUP BY date
			)
			SELECT
				owner,
				repository,
				A.date,
				merged,
				IFNULL(failure, 0) as failure ,
				IFNULL(failure, 0) / merged as failure_rate,
				IFNULL(MTTR, 0) as MTTR
			FROM A
			LEFT JOIN B
			ON A.date = B.date
			ORDER BY date DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
