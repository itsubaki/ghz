package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func PullReqsMeta(projectID, datasetName, tableName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				DATE_ADD(DATE(created_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(created_at), INTERVAL -0 DAY)) +1 DAY) as week,
				count(owner) / 7 as merged_per_day,
				AVG(TIMESTAMP_DIFF(merged_at, created_at,MINUTE)) as duration_avg,
				STDDEV(TIMESTAMP_DIFF(merged_at, created_at,MINUTE)) as duration_stddev
			FROM %v
			WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
			GROUP BY owner, repository, week
			ORDER BY week DESC
			LIMIT 1000
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, tableName),
		),
	}
}
