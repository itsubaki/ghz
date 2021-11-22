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
