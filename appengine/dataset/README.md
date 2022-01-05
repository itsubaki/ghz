# Query Example

## Lead Time

```sql
WITH A AS (
  SELECT
    owner,
    repository,
    workflow_name,
    DATE(completed_at) as date,
    PERCENTILE_CONT(lead_time, 0.5) OVER(partition by DATE(completed_at)) as lead_time
  FROM `$PROJECT_ID.vercel_next_js._leadtime_via_pullreqs`
)
SELECT
    owner,
    repository,
    workflow_name,
    date,
    MAX(lead_time) as lead_time
FROM A
GROUP BY owner, repository, workflow_name, date
ORDER BY date DESC
```

```json
[
  {
    "owner": "vercel",
    "repository": "next.js",
    "workflow_name": "Build, test, and deploy",
    "date": "2022-01-04",
    "lead_time": "79.0"
  },
  {
    "owner": "vercel",
    "repository": "next.js",
    "workflow_name": "Build, test, and deploy",
    "date": "2022-01-03",
    "lead_time": "16808.5"
  }
]
```

## MTTR

```sql
WITH A AS (
  SELECT
    owner,
    repository,
    DATE(merged_at) as date,
    PERCENTILE_CONT(TTR, 0.5) OVER(partition by DATE(merged_at)) as MTTR
  FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs`
)
SELECT
  owner,
  repository,
  date,
  MAX(MTTR) as MTTR
FROM A
GROUP BY owner, repository, date
ORDER BY date DESC
```

```json
[
  {
    "owner": "itsubaki",
    "repository": "ghz",
    "date": "2021-12-08",
    "MTTR": "60.0"
  }
]
```

## Failure Rate

```sql
WITH A AS (
  SELECT
    owner,
    repository,
    DATE(pushed_at) as date,
    COUNT(*) as failure
  FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pushed`
  GROUP BY date, owner, repository
), B AS (
  SELECT
    DATE(created_at) as date,
    COUNT(*) as pushed
  FROM `$PROJECT_ID.itsubaki_ghz.events_push`
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
```

```json
[
  {
    "date": "2021-12-08",
    "pushed": "13",
    "failure": "1",
    "failure_rate": "0.07692307692307693"
  },
  {
    "date": "2021-12-24",
    "pushed": "3",
    "failure": "1",
    "failure_rate": "0.3333333333333333"
  }
]
```

## Insert Incident

```sql
INSERT INTO `$PROJECT_ID.itsubaki_q.incidents` (
  owner,
  repository,
  id,
  description,
  sha,
  resolved_at
)
VALUES (
  'itsubaki',
  'q',
  '1',
  '[TEST] Incident via PullRequest',
  '7b2619e89065d96e683d70a72512e2883c1a2cf6',
  '2021-07-30 13:04:37 UTC'
)
```

```sql
INSERT INTO `$PROJECT_ID.itsubaki_q.incidents` (
  owner,
  repository,
  id,
  description,
  sha,
  resolved_at
)
VALUES (
  'itsubaki',
  'q',
  '1',
  '[TEST] Incident via Commit',
  'ad79208ce9ad1fce87b298ae28c6c518dc2a0486',
  '2021-12-26 15:31:05 UTC'
)
```

## JSON_EXTRACT

```sql
CREATE TEMP FUNCTION JSON2ARRAY(json STRING)
RETURNS ARRAY<STRING>
LANGUAGE js AS """
  return JSON.parse(json).map(x=>JSON.stringify(x));
""";

WITH A AS (
SELECT
  owner,
  repository,
  id,
  login,
  type,
  created_at,
  JSON_EXTRACT_SCALAR(raw_payload,'$.head') as head_sha,
  JSON2ARRAY(JSON_EXTRACT(raw_payload,'$.commits')) as commits
FROM `$PROJECT_ID.itsubaki_q.events`
WHERE type = "PushEvent"
)

SELECT
  owner,
  repository,
  id,
  login,
  type,
  created_at,
  head_sha,
  JSON_EXTRACT_SCALAR(commit,'$.sha') as sha
FROM A, UNNEST(commits) AS commit
```

```json
[
  {
    "owner": "itsubaki",
    "repository": "q",
    "id": "19394201253",
    "login": "itsubaki",
    "type": "PushEvent",
    "created_at": "2021-12-18 02:42:02 UTC",
    "head_sha": "2c76bb8c5e18ec7652a5205b294fec46f888fd52",
    "sha": "2c76bb8c5e18ec7652a5205b294fec46f888fd52"
  },
  {
    "owner": "itsubaki",
    "repository": "q",
    "id": "19394167370",
    "login": "itsubaki",
    "type": "PushEvent",
    "created_at": "2021-12-18 02:33:12 UTC",
    "head_sha": "667231cbfd88c9162f986e6021dd6303151230a4",
    "sha": "667231cbfd88c9162f986e6021dd6303151230a4"
  },
  {
    "owner": "itsubaki",
    "repository": "q",
    "id": "19221665438",
    "login": "itsubaki",
    "type": "PushEvent",
    "created_at": "2021-12-07 14:59:05 UTC",
    "head_sha": "42b43a568b29448e0bc60fecf8f94aa3df1c2798",
    "sha": "d00f69dcfa519148b769c2e2c9d7495e2a16b731"
  },
  {
    "owner": "itsubaki",
    "repository": "q",
    "id": "19221665438",
    "login": "itsubaki",
    "type": "PushEvent",
    "created_at": "2021-12-07 14:59:05 UTC",
    "head_sha": "42b43a568b29448e0bc60fecf8f94aa3df1c2798",
    "sha": "42b43a568b29448e0bc60fecf8f94aa3df1c2798"
  }
]
```

## Weekly

```sql
SELECT
  DATE_ADD(DATE(date), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(date), INTERVAL -0 DAY)) +1 DAY) as week,
  COUNT(merged) as merged
FROM `$PROJECT_ID.itsubaki_q._pullreqs`
GROUP BY week
```

```json
[
  {
    "week": "2021-07-25",
    "merged": "1"
  },
  {
    "week": "2018-07-29",
    "merged": "1"
  }
]
```
