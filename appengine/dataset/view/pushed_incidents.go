package view

import (
	"fmt"

	"cloud.google.com/go/bigquery"
	"github.com/itsubaki/ghz/appengine/dataset"
)

func PushedIncidentsMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed_incidents",
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
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.IncidentsMeta.Name),
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.EventsPushMeta.Name),
		),
	}
}

func PushedMTTRMeta(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed_mttr",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS(
				SELECT
					owner,
					repository,
					DATE(pushed_at) as date,
					PERCENTILE_CONT(TTR, 0.5) OVER(partition by DATE(pushed_at)) as MTTR
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
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PushedIncidentsMeta(id, dsn).Name),
		),
	}
}

func PushedFailureRate(id, dsn string) bigquery.TableMetadata {
	return bigquery.TableMetadata{
		Name: "_pushed_failure_rate",
		ViewQuery: fmt.Sprintf(
			`
			WITH A AS (
				SELECT
					owner,
					repository,
					DATE(pushed_at) as date,
					COUNT(*) as failure
				FROM %v
				GROUP BY date, owner, repository
			), B AS (
				SELECT
					DATE(created_at) as date,
					COUNT(*) as pushed
				FROM %v
				GROUP BY date
			)
			SELECT
				A.owner,
				A.repository,
				A.date,
				B.pushed,
				A.failure,
				A.failure / B.pushed as failure_rate
			FROM A
			INNER JOIN B
			ON A.date = B.date
			`,
			fmt.Sprintf("`%v.%v.%v`", id, dsn, PushedIncidentsMeta(id, dsn).Name),
			fmt.Sprintf("`%v.%v.%v`", id, dsn, dataset.EventsPushMeta.Name),
		),
	}
}
