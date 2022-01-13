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
			SELECT
				A.owner,
				A.repository,
				A.description,
				B.message,
				B.head_sha,
				B.sha,
				B.created_at as pushed_at,
				A.resolved_at,
				TIMESTAMP_DIFF(A.resolved_at, B.created_at, MINUTE) as TTR
			FROM %v as A
			INNER JOIN %v as B
			ON A.sha = B.sha
			`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.IncidentsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, dataset.EventsPushMeta.Name),
		),
	}
}
