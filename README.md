# prstats

* One of the indicators of productivity
* [DORA's Four Keys](https://github.com/GoogleCloudPlatform/fourkeys)

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```shell
$ prstats runslist --owner itsubaki --repo q  > out.json
$ prstats analyze --path out.json --format csv | column -t -s, | less -S
workflow_ID   name    start                           end                             run_per_day   failure_rate   duration_avg(m)
...
5841880       tests   2021-02-14 00:00:00 +0000 UTC   2021-02-21 00:00:00 +0000 UTC   0.28          0              0.59
5841880       tests   2021-02-21 00:00:00 +0000 UTC   2021-02-28 00:00:00 +0000 UTC   0.57          0              0.72
5841880       tests   2021-02-28 00:00:00 +0000 UTC   2021-03-07 00:00:00 +0000 UTC   2.14          0.066          2.46
5841880       tests   2021-03-07 00:00:00 +0000 UTC   2021-03-14 00:00:00 +0000 UTC   3.71          0.038          2.53
5841880       tests   2021-03-14 00:00:00 +0000 UTC   2021-03-21 00:00:00 +0000 UTC   0.42          0              0.97
...
```

```shell
$ prstats runslist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
workflow_ID   name    number   run_ID       conclusion   status      created_at                      updated_at                      duration(m)
6067686       tests   77       1429354728   success      completed   2021-11-06 15:04:28 +0000 UTC   2021-11-06 15:07:13 +0000 UTC   2.75
6067686       tests   76       1245764204   success      completed   2021-09-17 14:03:30 +0000 UTC   2021-09-17 14:06:06 +0000 UTC   2.6
6067686       tests   75       1224424786   success      completed   2021-09-11 13:31:18 +0000 UTC   2021-09-11 13:33:57 +0000 UTC   2.65
6067686       tests   74       1224410044   failure      completed   2021-09-11 13:22:05 +0000 UTC   2021-09-11 13:24:46 +0000 UTC   2.68
6067686       tests   73       1224351644   success      completed   2021-09-11 12:48:37 +0000 UTC   2021-09-11 12:50:55 +0000 UTC   2.3
6067686       tests   72       1224334415   success      completed   2021-09-11 12:37:14 +0000 UTC   2021-09-11 12:39:36 +0000 UTC   2.36
6067686       tests   71       1224320650   success      completed   2021-09-11 12:28:25 +0000 UTC   2021-09-11 12:31:13 +0000 UTC   2.8
6067686       tests   70       1224306965   success      completed   2021-09-11 12:20:32 +0000 UTC   2021-09-11 12:21:38 +0000 UTC   1.1
...
```

```shell
$ prstats prlist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
id          title                        created_at                      merged_at                       durationm)   
545593516   gorm v2                      2020-12-25 13:37:55 +0000 UTC   2020-12-26 08:58:26 +0000 UTC   1160.5167       
473905785   Feature/godog v0.10.0        2020-08-26 13:25:52 +0000 UTC   2020-08-26 13:26:12 +0000 UTC   0.3333          
425425099   Rename repository            2020-05-30 06:47:29 +0000 UTC   2020-05-30 06:47:38 +0000 UTC   0.1500          
306753867   Refactor repository          2019-08-13 05:44:48 +0000 UTC   2019-08-13 05:50:56 +0000 UTC   6.1333          
282046942   Feature/multitenant          2019-05-24 14:47:50 +0000 UTC   2019-05-24 14:50:02 +0000 UTC   2.2000          
274962862   Applied Clean Architecture   2019-05-01 06:00:21 +0000 UTC   2019-05-01 06:00:30 +0000 UTC   0.1500            
```
