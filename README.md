# prstats

* One of the indicators of productivity
* [DORA Four Keys](https://github.com/GoogleCloudPlatform/fourkeys)

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```shell
$ prstats runslist --owner itsubaki --repo q > out.json
$ prstats analyze  --path out.json --format csv | column -t -s, | less -S
workflow_ID   name    start        end          runs_per_day   failure_rate   duration_avg(m)
...
5841880       tests   2021-02-14   2021-02-21   0.28           0              0.59
5841880       tests   2021-02-21   2021-02-28   0.57           0              0.72
5841880       tests   2021-02-28   2021-03-07   2.14           0.066          2.46
5841880       tests   2021-03-07   2021-03-14   3.71           0.038          2.53
5841880       tests   2021-03-14   2021-03-21   0.42           0              0.97
...
```


```shell
$ prstats runslist --owner itsubaki --repo q > out.json
$ prstats jobslist --owner itsubaki --repo q --path out.json --format csv | column -t -s, | less -S
workflow_name   run_id       run_number   job_id       job_name               conclusion       status      started_at            completed_at          duration(minutes)
tests           1429332756   109          4126276932   test (ubuntu-latest)   success          completed   2021-11-06 14:52:00   2021-11-06 14:52:53   0.88
tests           1299800366   108          3780280583   test (ubuntu-latest)   success          completed   2021-10-03 05:11:05   2021-10-03 05:11:43   0.63
tests           1200937783   107          3513358909   test (ubuntu-latest)   success          completed   2021-09-04 12:41:10   2021-09-04 12:41:53   0.72
tests           1174489906   106          3443854949   test (ubuntu-latest)   success          completed   2021-08-27 14:00:47   2021-08-27 14:01:36   0.82
tests           1169027869   105          3428886296   test (ubuntu-latest)   success          completed   2021-08-26 03:51:14   2021-08-26 03:51:56   0.7
tests           1166653588   104          3422345982   test (ubuntu-latest)   success          completed   2021-08-25 13:20:52   2021-08-25 13:21:54   1.03
...
```

```shell
$ prstats runslist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
workflow_ID   name    number   run_ID       conclusion   status      created_at            updated_at            duration(m)
6067686       tests   77       1429354728   success      completed   2021-11-06 15:04:28   2021-11-06 15:07:13   2.75
6067686       tests   76       1245764204   success      completed   2021-09-17 14:03:30   2021-09-17 14:06:06   2.6
6067686       tests   75       1224424786   success      completed   2021-09-11 13:31:18   2021-09-11 13:33:57   2.65
6067686       tests   74       1224410044   failure      completed   2021-09-11 13:22:05   2021-09-11 13:24:46   2.68
6067686       tests   73       1224351644   success      completed   2021-09-11 12:48:37   2021-09-11 12:50:55   2.3
6067686       tests   72       1224334415   success      completed   2021-09-11 12:37:14   2021-09-11 12:39:36   2.37
6067686       tests   71       1224320650   success      completed   2021-09-11 12:28:25   2021-09-11 12:31:13   2.8
6067686       tests   70       1224306965   success      completed   2021-09-11 12:20:32   2021-09-11 12:21:38   1.1
6067686       tests   69       1224295793   success      completed   2021-09-11 12:13:52   2021-09-11 12:14:59   1.11
6067686       tests   68       1209535428   success      completed   2021-09-07 13:07:07   2021-09-07 13:08:17   1.17
6067686       tests   67       1209364891   success      completed   2021-09-07 12:14:19   2021-09-07 12:15:44   1.41
6067686       tests   66       1185775378   success      completed   2021-08-31 10:36:07   2021-08-31 10:37:12   1.08
6067686       tests   65       1185771230   success      completed   2021-08-31 10:34:51   2021-08-31 10:35:52   1.01
6067686       tests   64       1116595181   success      completed   2021-08-10 12:30:51   2021-08-10 12:32:05   1.23
6067686       tests   63       1092702433   success      completed   2021-08-03 04:48:42   2021-08-03 04:49:49   1.11

...
```

```shell
$ prstats prlist --owner itsubaki --repo mackerel-server-go --format csv | column -t -s, | less -S
id          title                        created_at            merged_at             duration(m)   
545593516   gorm v2                      2020-12-25 13:37:55   2020-12-26 08:58:26   1160.5167     
473905785   Feature/godog v0.10.0        2020-08-26 13:25:52   2020-08-26 13:26:12   0.3333        
425425099   Rename repository            2020-05-30 06:47:29   2020-05-30 06:47:38   0.1500        
306753867   Refactor repository          2019-08-13 05:44:48   2019-08-13 05:50:56   6.1333        
282046942   Feature/multitenant          2019-05-24 14:47:50   2019-05-24 14:50:02   2.2000        
274962862   Applied Clean Architecture   2019-05-01 06:00:21   2019-05-01 06:00:30   0.1500           
```
