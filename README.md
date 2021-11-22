# prstats
CLI for Github PR stats

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```shell
$ prstats --owner itsubaki --repo q | jq .
{
  "owner": "itsubaki",
  "repo": "q",
  "range": {
    "beg": "2021-07-27T13:40:39Z",
    "end": "2018-07-29T12:41:02Z"
  },
  "pr": {
    "count_per_day": 0.003656307129798903,
    "count": 4,
    "days": 1094
  },
  "merged": {
    "count_per_day": 0.0018248175182481751,
    "count": 2,
    "days": 1096,
    "hours_per_count": 35,
    "total_hours": 70
  },
  "workflow": [
    {
      "id": 5841880,
      "name": "tests",
      "failure_rate": 0.03669724770642202,
      "count": 109,
      "success": 105,
      "failure": 4,
      "skipped": 0,
      "cancelled": 0
    }
  ]
}
```


```shell
$ prstats --owner itsubaki --repo mackerel-server-go list --format csv | column -t -s, | less -S
id          title                        created_at                      merged_at                       lead_time(hours)   
545593516   gorm v2                      2020-12-25 13:37:55 +0000 UTC   2020-12-26 08:58:26 +0000 UTC   19.3419            
473905785   Feature/godog v0.10.0        2020-08-26 13:25:52 +0000 UTC   2020-08-26 13:26:12 +0000 UTC   0.0056             
425425099   Rename repository            2020-05-30 06:47:29 +0000 UTC   2020-05-30 06:47:38 +0000 UTC   0.0025             
306753867   Refactor repository          2019-08-13 05:44:48 +0000 UTC   2019-08-13 05:50:56 +0000 UTC   0.1022             
282046942   Feature/multitenant          2019-05-24 14:47:50 +0000 UTC   2019-05-24 14:50:02 +0000 UTC   0.0367             
274962862   Applied Clean Architecture   2019-05-01 06:00:21 +0000 UTC   2019-05-01 06:00:30 +0000 UTC   0.0025             
```