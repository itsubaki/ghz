package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PullReqsTTRMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_ttr",
		ViewQuery: fmt.Sprintf(
			`
			SELECT
				A.owner,
				A.repository,
				A.description,
				B.id,
				B.number,
				B.login,
				B.title,
				B.message,
				B.merge_commit_sha,
				A.sha,
				B.merged_at,
				A.resolved_at,
				TIMESTAMP_DIFF(A.resolved_at, B.merged_at, MINUTE) as TTR
			FROM %v as A
			INNER JOIN %v as B
			ON A.sha = B.sha
			`,
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.IncidentsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PullReqsMeta(id, dsn).Name),
		),
	}
}

func PullReqsTTRMedianMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_ttr_median",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS(
				SELECT
					owner,
					repository,
					DATE(merged_at) as date,
					PERCENTILE_CONT(TTR, 0.5) OVER(partition by DATE(merged_at)) as MTTR
				FROM %v
			)
			SELECT
				owner,
				repository,
				date,
				MAX(MTTR) as MTTR
			FROM A
			GROUP BY owner, repository, date
			`,
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PullReqsTTRMeta(id, dsn).Name),
		),
	}
}

func PullReqsFailureRate(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pullreqs_failure_rate",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					owner,
					repository,
					DATE(merged_at) as date,
					COUNT(*) as failure
				FROM %v
				GROUP BY date, owner, repository
			), B AS (
				SELECT
					DATE(merged_at) as date,
					COUNT(*) as merged
				FROM %v
				WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
				GROUP BY date
			)
			SELECT
				A.owner,
				A.repository,
				A.date,
				B.merged,
				A.failure,
				A.failure / B.merged as failure_rate
			FROM A
			INNER JOIN B
			ON A.date = B.date
			`,
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PullReqsTTRMeta(id, dsn).Name),
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.PullReqsMeta.Name),
		),
	}
}
