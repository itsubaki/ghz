package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PullReqsLeadTimeMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_leadtime",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.workflow_id,
				A.workflow_name,
				B.id,
				B.number,
				B.login,
				B.title,
				B.message,
				B.merge_commit_sha,
				B.sha,
				B.committed_at,
				A.updated_at as completed_at,
				TIMESTAMP_DIFF(A.updated_at, B.committed_at, MINUTE) as lead_time
			FROM %v as A
			INNER JOIN %v as B
			ON A.head_sha = B.merge_commit_sha
			WHERE A.conclusion = "success"
			`,
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.WorkflowRunsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PullReqsMeta(id, dsn).Name),
		),
	}
}

func PullReqsLeadTimeMedianMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_leadtime_median",
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
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PullReqsLeadTimeMeta(id, dsn).Name),
		),
	}
}
