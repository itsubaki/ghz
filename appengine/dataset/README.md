# Query Example

## Insert

```sql
INSERT INTO `PROJECT_ID.DATASET_NAME.incidents` (
  owner,
  repository,
  id,
  description,
  sha,
  created_at,
  resolved_at
)
VALUES (
  'itsubaki',
  'q',
  '1',
  'TEST',
  '2c76bb8c5e18ec7652a5205b294fec46f888fd52',
  '2021-12-18 02:42:02 UTC',
  '2021-12-24 02:42:02 UTC'
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
FROM `PROJECT_ID.DATASET_NAME.events`
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
