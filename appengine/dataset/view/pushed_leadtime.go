package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PushedLeadTimeMeta(dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed_leadtime",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.workflow_id,
				A.workflow_name,
				B.login,
				B.message,
				B.head_sha,
				B.sha,
				B.committed_at,
				A.updated_at as completed_at,
				TIMESTAMP_DIFF(A.updated_at, B.committed_at, MINUTE) as lead_time
			FROM %v as A
			INNER JOIN %v as B
			ON A.head_sha = B.head_sha
			AND A.conclusion = "success"
			`,
			fmt.Sprintf("`%v.%v.%v`", dataset.ProjectID, dsn, dataset.WorkflowRunsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", dataset.ProjectID, dsn, PushedMeta(dsn).Name),
		),
	}
}

func PushedLeadTimeMedianMeta(dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed_leadtime_median",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					owner,
					repository,
					workflow_name,
					DATE(completed_at) as date,
					PERCENTILE_CONT(lead_time, 0.5) OVER(partition by DATE(completed_at)) as lead_time
				FROM %v
			)
			SELECT
				owner,
				repository,
				workflow_name,
				date,
				MAX(lead_time) as lead_time
			FROM A
			GROUP BY owner, repository, workflow_name, date
			`,
			fmt.Sprintf("`%v.%v.%v`", dataset.ProjectID, dsn, PushedLeadTimeMeta(dsn).Name),
		),
	}
}
