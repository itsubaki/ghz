package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghstats/appengine/dataset"
)

func IncidentsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_incidents",
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
					count(created_at) as failure,
					DATE_ADD(DATE(created_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(created_at), INTERVAL -0 DAY)) +1 DAY) as week
				FROM %v
				GROUP BY week
			)
			SELECT
				owner,
				repository,
				A.week,
				commits,
				IFNULL(failure, 0) as failure ,
				IFNULL(failure,0)/commits as failure_rate
			FROM A
			LEFT JOIN B
			ON A.week = B.week
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.CommitsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
		),
	}
}
