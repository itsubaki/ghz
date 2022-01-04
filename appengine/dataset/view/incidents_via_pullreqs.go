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
					COUNT(merged_at) as merged,
					DATE(merged_at) as date
				FROM %v
				WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
				GROUP BY owner, repository, date
			), B AS (
				SELECT
					B.sha,
					B.date,
					A.merged_at
				FROM %v as A
				INNER JOIN %v as B
				ON A.id = B.id
				WHERE A.state = "closed" AND A.merged_at != "0001-01-01 00:00:00 UTC"
			), C AS (
				SELECT
					Date(B.merged_at) as date,
					COUNT(B.merged_at) as failure,
				FROM %v as A
				INNER JOIN B
				ON A.sha = B.sha
				GROUP BY date
			), D AS (
				SELECT
					Date(B.merged_at) as date,
                	PERCENTILE_CONT(TIMESTAMP_DIFF(A.resolved_at, B.merged_at, MINUTE),0.5) OVER(partition by Date(B.merged_at)) as MTTR
				FROM %v as A
				INNER JOIN B
				ON A.sha = B.sha
			), E AS (
				SELECT
					C.date,
					C.failure,
					D.MTTR,
				FROM C
				INNER JOIN D
				ON C.date = D.date
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
			LEFT JOIN E
			ON A.date = E.date
			ORDER BY date DESC
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqCommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
