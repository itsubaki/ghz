# prstats
CLI for Github PR stats

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```
$ prstats --org itsubaki --repo q | jq .
{
  "org": "itsubaki",
  "repo": "q",
  "range": {
    "beg": "2021-07-27T13:40:39Z",
    "end": "2018-07-29T12:41:02Z"
  },
  "pr": {
    "count": 4,
    "days": 1094,
    "rate": 0.003656307129798903
  },
  "lifetime": {
    "average_hours": 35,
    "total_hours": 70,
    "count": 2
  }
}
```
