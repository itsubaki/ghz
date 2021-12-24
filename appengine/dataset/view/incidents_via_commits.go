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
					DATE(date) as date,
					count(date) as commits
				FROM %v
				GROUP BY owner, repository, date
			), B AS (
				SELECT
					DATE(A.date) as date,
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
				commits,
				IFNULL(failure, 0) as failure,
				IFNULL(failure, 0) / commits as failure_rate,
				IFNULL(MTTR, 0) as MTTR
			FROM A
			LEFT JOIN B
			ON A.date = B.date
			ORDER BY date DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
