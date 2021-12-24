package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
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
					DATE_ADD(DATE(merged_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(merged_at), INTERVAL -0 DAY)) +1 DAY) as week
				FROM %v
				WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
				GROUP BY owner, repository, week
			), B AS (
				SELECT
					avg(TIMESTAMP_DIFF(B.resolved_at, A.date, MINUTE)) as MTTR,
					count(A.date) as failure,
					DATE_ADD(DATE(A.date), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(A.date), INTERVAL -0 DAY)) +1 DAY) as week
				FROM %v as A
				INNER JOIN %v as B
				ON A.sha = B.sha
				GROUP BY week
			)
			SELECT
				owner,
				repository,
				A.week,
				merged,
				IFNULL(failure, 0) as failure ,
				IFNULL(failure, 0) / merged as failure_rate,
				IFNULL(MTTR, 0) as MTTR
			FROM A
			LEFT JOIN B
			ON A.week = B.week
			ORDER BY week DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
