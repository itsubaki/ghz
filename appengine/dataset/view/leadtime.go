package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
)

func LeadTimeMeta(projectID, datasetName, tableName string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_leadtime",
		ViewQuery: fmt.Sprintf(
			`SELECT * FROM %v`,
			fmt.Sprintf("`%v.%v.%v`", projectID, datasetName, tableName),
		),
	}
}
