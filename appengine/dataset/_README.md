```sql
SELECT
 DATE_ADD(DATE(date), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(date), INTERVAL -0 DAY)) +1 DAY) as week,
 count(sha) / 5 as commit_per_day,
FROM `$PROJECT_ID.$DATASET.commits`
GROUP BY week

[
  {
    "week": "2021-12-05",
    "commit_per_day": "0.6"
  },
  {
    "week": "2021-10-31",
    "commit_per_day": "0.2"
  },
...
]
```

```sql
SELECT
 DATE_ADD(DATE(merged_at), INTERVAL - EXTRACT(DAYOFWEEK FROM DATE_ADD(DATE(merged_at), INTERVAL -0 DAY)) +1 DAY) as week,
 count(id) / 5 as merged_per_day,
FROM `$PROJECT_ID.$DATASET.pullreqs`
WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
GROUP BY week
```
