package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func PullReqsCommitsMeta(projectID, datasetName, tableName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_commits",
		ViewQuery: fmt.Sprintf(
			"SELECT * FROM %v",
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, tableName),
		),
	}
}
