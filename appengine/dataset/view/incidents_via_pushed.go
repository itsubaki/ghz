package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func IncidentsPushedMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_incidents_via_pushed",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					owner,
					repository,
					DATE(created_at) as date,
					COUNT(DISTINCT(head_sha)) as pushed
				FROM %v
				GROUP BY owner, repository, date
			), B AS (
				SELECT
					DATE(A.created_at) as date,
					COUNT(A.created_at) as failure,
					AVG(TIMESTAMP_DIFF(B.resolved_at, A.created_at, MINUTE)) as MTTR
				FROM %v as A
				INNER JOIN %v as B
				ON A.sha = B.sha
				GROUP BY date
			)
			SELECT
				owner,
				repository,
				A.date,
				pushed,
				IFNULL(failure, 0) as failure,
				IFNULL(failure, 0) / pushed as failure_rate,
				IFNULL(MTTR, 0) as MTTR
			FROM A
			LEFT JOIN B
			ON A.date = B.date
			ORDER BY date DESC
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.EventsPushMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.EventsPushMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
