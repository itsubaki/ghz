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
...
workflow_ID   name    start                           end                             run_per_day           failure_rate           duration_avg(min)
5841880       tests   2021-02-14 00:00:00 +0000 UTC   2021-02-21 00:00:00 +0000 UTC   0.2857142857142857    0                      0.5916666666666667
5841880       tests   2021-02-21 00:00:00 +0000 UTC   2021-02-28 00:00:00 +0000 UTC   0.5714285714285714    0                      0.7291666666666666
5841880       tests   2021-02-28 00:00:00 +0000 UTC   2021-03-07 00:00:00 +0000 UTC   2.142857142857143     0.06666666666666667    2.4655555555555555
5841880       tests   2021-03-07 00:00:00 +0000 UTC   2021-03-14 00:00:00 +0000 UTC   3.7142857142857144    0.038461538461538464   2.532051282051282
5841880       tests   2021-03-14 00:00:00 +0000 UTC   2021-03-21 00:00:00 +0000 UTC   0.42857142857142855   0                      0.9777777777777779
...
```

```shell
$ prstats runslist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
workflow_ID   name    number   run_ID       conclusion   status      created_at                      updated_at                      duration(hours)
6067686       tests   77       1429354728   success      completed   2021-11-06 15:04:28 +0000 UTC   2021-11-06 15:07:13 +0000 UTC   0.04583333333333333
6067686       tests   76       1245764204   success      completed   2021-09-17 14:03:30 +0000 UTC   2021-09-17 14:06:06 +0000 UTC   0.043333333333333335
6067686       tests   75       1224424786   success      completed   2021-09-11 13:31:18 +0000 UTC   2021-09-11 13:33:57 +0000 UTC   0.04416666666666667
6067686       tests   74       1224410044   failure      completed   2021-09-11 13:22:05 +0000 UTC   2021-09-11 13:24:46 +0000 UTC   0.04472222222222222
6067686       tests   73       1224351644   success      completed   2021-09-11 12:48:37 +0000 UTC   2021-09-11 12:50:55 +0000 UTC   0.03833333333333333
6067686       tests   72       1224334415   success      completed   2021-09-11 12:37:14 +0000 UTC   2021-09-11 12:39:36 +0000 UTC   0.03944444444444444
6067686       tests   71       1224320650   success      completed   2021-09-11 12:28:25 +0000 UTC   2021-09-11 12:31:13 +0000 UTC   0.04666666666666667
6067686       tests   70       1224306965   success      completed   2021-09-11 12:20:32 +0000 UTC   2021-09-11 12:21:38 +0000 UTC   0.018333333333333333
...
```

```shell
$ prstats prlist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
id          title                        created_at                      merged_at                       duration(hours)   
545593516   gorm v2                      2020-12-25 13:37:55 +0000 UTC   2020-12-26 08:58:26 +0000 UTC   19.3419            
473905785   Feature/godog v0.10.0        2020-08-26 13:25:52 +0000 UTC   2020-08-26 13:26:12 +0000 UTC   0.0056             
425425099   Rename repository            2020-05-30 06:47:29 +0000 UTC   2020-05-30 06:47:38 +0000 UTC   0.0025             
306753867   Refactor repository          2019-08-13 05:44:48 +0000 UTC   2019-08-13 05:50:56 +0000 UTC   0.1022             
282046942   Feature/multitenant          2019-05-24 14:47:50 +0000 UTC   2019-05-24 14:50:02 +0000 UTC   0.0367             
274962862   Applied Clean Architecture   2019-05-01 06:00:21 +0000 UTC   2019-05-01 06:00:30 +0000 UTC   0.0025             
...
```
