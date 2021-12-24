package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func IncidentsCommitsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_incidents_via_commits",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					owner,
					repository,
					count(date) as commits,
					DATE_ADD(DATE(date), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(date), INTERVAL -0 DAY)) +1 DAY) as week
				FROM %v
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
				commits,
				IFNULL(failure, 0) as failure ,
				IFNULL(failure, 0) / commits as failure_rate,
				IFNULL(MTTR, 0) as MTTR
			FROM A
			LEFT JOIN B
			ON A.week = B.week
			ORDER BY week DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
