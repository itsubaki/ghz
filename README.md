# prstats

* One of the indicators of productivity

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```shell
$ prstats --owner itsubaki --repo q --days 365 | jq .
{
  "owner": "itsubaki",
  "repo": "q",
  "range": {
    "beg": "2020-11-24T10:30:28.70222+09:00",
    "end": "2021-11-24T10:30:28.70222+09:00",
    "days": 365
  },
  "created": {
    "count_per_day": 0.0027,
    "count": 1
  },
  "merged": {
    "count_per_created": 1,
    "count_per_day": 0.0027,
    "hours_per_count": 70.39,
    "total_hours": 70.39,
    "count": 1
  },
  "workflow_runs": [
    {
      "id": 5841880,
      "name": "tests",
      "count_per_day": 0.29,
      "failure_rate": 0.036,
      "duration_hours": 6.51,
      "success": 105,
      "failure": 4,
      "skipped": 0,
      "cancelled": 0,
      "action_required": 0,
      "count": 109
    }
  ]
}
```


```shell
$ prstats --owner itsubaki --repo mackerel-server-go --format csv list | column -t -s, | less -S
id          title                        created_at                      merged_at                       lead_time(hours)   
545593516   gorm v2                      2020-12-25 13:37:55 +0000 UTC   2020-12-26 08:58:26 +0000 UTC   19.3419            
473905785   Feature/godog v0.10.0        2020-08-26 13:25:52 +0000 UTC   2020-08-26 13:26:12 +0000 UTC   0.0056             
425425099   Rename repository            2020-05-30 06:47:29 +0000 UTC   2020-05-30 06:47:38 +0000 UTC   0.0025             
306753867   Refactor repository          2019-08-13 05:44:48 +0000 UTC   2019-08-13 05:50:56 +0000 UTC   0.1022             
282046942   Feature/multitenant          2019-05-24 14:47:50 +0000 UTC   2019-05-24 14:50:02 +0000 UTC   0.0367             
274962862   Applied Clean Architecture   2019-05-01 06:00:21 +0000 UTC   2019-05-01 06:00:30 +0000 UTC   0.0025             
```
