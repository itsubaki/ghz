# prstats
CLI for Github PR stats

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```
$ prstats pr --org itsubaki --repo q
ID: 697873132, title: Add lexer for OpenQASM, created: 2021-07-27 13:40:39 +0000 UTC, merged: 2021-07-30 12:04:37 +0000 UTC, closed: 2021-07-30 12:04:37 +0000 UTC
ID: 444592694, title: Configure Sider, created: 2020-07-06 07:30:06 +0000 UTC, merged: <nil>, closed: 2020-07-07 00:43:12 +0000 UTC
ID: 270932381, title: Added documentation and expanded coverage, created: 2019-04-16 14:21:28 +0000 UTC, merged: <nil>, closed: 2021-02-25 12:08:43 +0000 UTC
ID: 204640368, title: Add simulator, created: 2018-07-29 12:41:02 +0000 UTC, merged: 2018-07-29 12:42:47 +0000 UTC, closed: 2018-07-29 12:42:47 +0000 UTC
```
