package dataset

import "cloud.google.com/go/bigquery"

type Tag struct {
	Owner      string `bigquery:"owner"`
	Repository string `bigquery:"repository"`
	Name       string `bigquery:"name"`
	SHA        string `bigquery:"tag"`
}

var TagsMeta = bigquery.TableMetadata{
	Name: "tags",
	Schema: bigquery.Schema{
		{Name: "owner", Type: bigquery.StringFieldType, Required: true},
		{Name: "repository", Type: bigquery.StringFieldType, Required: true},
		{Name: "name", Type: bigquery.StringFieldType, Required: true},
		{Name: "sha", Type: bigquery.StringFieldType, Required: true},
	},
}
