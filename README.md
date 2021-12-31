# ghz

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/ghz)](https://pkg.go.dev/github.com/itsubaki/ghz)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/ghz?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/ghz)

- One of the indicators of productivity
- [DORA's Four Keys](https://github.com/GoogleCloudPlatform/fourkeys)
  - Deployment Frequency: The number of deployments to production per day.
  - Lead Time for Changes: The difference between the commit creation time and the deployment completion time.
  - Time to Restore Services: The difference between the commit creation time that caused the failure and the resolution time.
  - Change Failure Rate: The Ratio of the number of deployments and the number of failures.

## Install

```shell
go install github.com/itsubaki/ghz@latest
```

## Example

```shell
$ ghz commits fetch --owner itsubaki -repo q
$ ghz commits list  --owner itsubaki -repo q --format csv | column -t -s, | less -S
sha                                        login      date                  message
42b43a568b29448e0bc60fecf8f94aa3df1c2798   itsubaki   2021-12-07 14:58:53   Merge branch 'main' of https://github.com/itsubaki/q into main
d00f69dcfa519148b769c2e2c9d7495e2a16b731   itsubaki   2021-12-07 14:58:40   Update test
09638c9e19af748d434fa9afbd3e48a4a5b74df1   itsubaki   2021-12-07 03:58:46   Refactor cmodexp2 gate
...
```

```shell
$ ghz pullreqs fetch --owner itsubaki -repo q
$ ghz pullreqs list  --owner itsubaki -repo q --format csv | column -t -s, | less -S
id          number   title                                       login        state    created_at            updated_at            merged_at             closed_at             merge_commit_sha
697873132   13       Add lexer for OpenQASM                      itsubaki     closed   2021-07-27 13:40:39   2021-07-30 12:05:17   2021-07-30 12:04:37   2021-07-30 12:04:37   6c6df53c0ee86e8e78e4313135ec11c0e6fa8764
444592694   11       Configure Sider                             sider[bot]   closed   2020-07-06 07:30:06   2021-11-22 00:43:05   null                  2020-07-07 00:43:12   87cb3d0103c04be2288e96a5acc08454c0b8788b
270932381   3        Added documentation and expanded coverage   axamon       closed   2019-04-16 14:21:28   2021-02-25 12:08:43   null                  2021-02-25 12:08:43   09ee1f68eba4eaad660cb48dd2237d0cb8c90495
204640368   2        Add simulator                               itsubaki     closed   2018-07-29 12:41:02   2021-07-24 11:27:42   2018-07-29 12:42:47   2018-07-29 12:42:47   dfdd80c575874dd4485007a5cda984e0b08a6ae8
```

```shell
$ ghz pullreqs commits fetch --owner itsubaki -repo q
$ ghz pullreqs commits list  --owner itsubaki -repo q --format csv | column -t -s, | less -S
pr_id       pr_number   sha                                        login               date                  message
697873132   13          7b2619e89065d96e683d70a72512e2883c1a2cf6   itsubaki            2021-07-30 12:02:49   Merge branch 'main' into openqasm
697873132   13          806cdd051a833de04ce1d3b721eff12004c64f41   itsubaki            2021-07-27 13:38:22   Add lexer for openqasm
...
```

```shell
$ ghz actions runs fetch --owner itsubaki -repo q
$ ghz actions runs list  --owner itsubaki -repo q --format csv | column -t -s, | less -S
workflow_id   workflow_name   run_id       run_number   status      conclusion   created_at            updated_at            head_commit.sha
5841880       tests           1549986886   111          completed   success      2021-12-07 14:59:05   2021-12-07 15:00:05   42b43a568b29448e0bc60fecf8f94aa3df1c2798
5841880       tests           1547804641   110          completed   success      2021-12-07 03:59:03   2021-12-07 04:00:04   09638c9e19af748d434fa9afbd3e48a4a5b74df1
...
5841880       tests           1082134326   91           completed   success      2021-07-30 12:04:39   2021-07-30 12:05:37   6c6df53c0ee86e8e78e4313135ec11c0e6fa8764
5841880       tests           1082128658   90           completed   success      2021-07-30 12:02:58   2021-07-30 12:03:54   7b2619e89065d96e683d70a72512e2883c1a2cf6
...
```

```shell
$ ghz actions jobs fetch --owner itsubaki -repo q
$ ghz actions jobs list  --owner itsubaki -repo q --format csv | column -t -s, | less -S
run_id       job_id       job_name               status           conclusion   started_at            completed_at
1549986886   4445591170   test (ubuntu-latest)   completed        success      2021-12-07 14:59:15   2021-12-07 15:00:01
1547804641   4439429997   test (ubuntu-latest)   completed        success      2021-12-07 03:59:10   2021-12-07 04:00:00
...
1082134326   3201660144   test (ubuntu-latest)   completed        success      2021-07-30 12:04:45   2021-07-30 12:05:33
1082128658   3201644714   test (ubuntu-latest)   completed        success      2021-07-30 12:03:06   2021-07-30 12:03:51
...
```

```shell
$ ghz pullreqs list --owner itsubaki -repo q | jq -r '[.id, .number, .user.login, .title, .state] | @csv' | tr -d '"' | column -t -s, | less -S
697873132  13  itsubaki    Add lexer for OpenQASM                     closed
444592694  11  sider[bot]  Configure Sider                            closed
270932381  3   axamon      Added documentation and expanded coverage  closed
204640368  2   itsubaki    Add simulator                              closed
```
