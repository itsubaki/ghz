package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PullReqsMeta(projectID, datasetName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				owner,
				repository,
				DATE(merged_at) as date,
				COUNT(owner) as merged,
				AVG(TIMESTAMP_DIFF(merged_at, created_at, MINUTE)) as duration_avg
			FROM %v
			WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
			GROUP BY owner, repository, date
			ORDER BY date DESC
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.PullReqsMeta.Name),
		),
	}
}
