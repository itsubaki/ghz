# prstats

* One of the indicators of productivity
* [DORA's Four Keys](https://github.com/GoogleCloudPlatform/fourkeys)

## Install

```shell
go install github.com/itsubaki/prstats@latest
```

## Example

```shell
$ prstats actions runs fetch --owner itsubaki --repo q
$ prstats actions runs list  --owner itsubaki --repo q --format csv | column -t -s, | less -S
workflow_id   workflow_name   run_id       run_number   status      conclusion   created_at            updated_at            duration(minutes)
5841880       tests           1429332756   109          completed   success      2021-11-06 14:51:50   2021-11-06 14:53:20   1.5
5841880       tests           1299800366   108          completed   success      2021-10-03 05:10:58   2021-10-03 05:11:45   0.7833333333333333
5841880       tests           1200937783   107          completed   success      2021-09-04 12:41:02   2021-09-04 12:41:56   0.9
5841880       tests           1174489906   106          completed   success      2021-08-27 14:00:40   2021-08-27 14:01:39   0.9833333333333333
5841880       tests           1169027869   105          completed   success      2021-08-26 03:51:07   2021-08-26 03:51:59   0.8666666666666667
5841880       tests           1166653588   104          completed   success      2021-08-25 13:20:43   2021-08-25 13:21:58   1.25
5841880       tests           1166541361   103          completed   success      2021-08-25 12:47:48   2021-08-25 12:50:36   2.8
5841880       tests           1166522425   102          completed   success      2021-08-25 12:42:05   2021-08-25 12:43:02   0.95
5841880       tests           1166480252   101          completed   success      2021-08-25 12:28:34   2021-08-25 12:29:38   1.0666666666666667
...

$ prstats actions runs analyze --owner itsubaki --repo q --format csv | column -t -s, | less -S
workflow_id   name    start        end          runs_per_day          failure_rate           duration_avg(minutes)   duration_var(minutes)
5841880       tests   2021-02-14   2021-02-21   0.2857142857142857    0                      0.5916666666666667      6.944444444444396e-05
5841880       tests   2021-02-21   2021-02-28   0.5714285714285714    0                      0.7291666666666666      0.07421875
5841880       tests   2021-02-28   2021-03-07   2.142857142857143     0.06666666666666667    2.4655555555555555      3.1072715226337455
5841880       tests   2021-03-07   2021-03-14   3.7142857142857144    0.038461538461538464   2.532051282051282       4.895141291154604
5841880       tests   2021-03-14   2021-03-21   0.42857142857142855   0                      0.9777777777777779      0.0008024691358024684
5841880       tests   2021-03-21   2021-03-28   0                     0                      0                       0
...
```

```shell
$ prstats actions jobs fetch --owner itsubaki --repo q
$ prstats actions jobs list  --owner itsubaki --repo q --format csv | column -t -s, | less -S
job_id       job_name               status           conclusion   started_at            completed_at          duration(minutes)
4126276932   test (ubuntu-latest)   completed        success      2021-11-06 14:52:00   2021-11-06 14:52:53   0.8833333333333333
3780280583   test (ubuntu-latest)   completed        success      2021-10-03 05:11:05   2021-10-03 05:11:43   0.6333333333333333
3513358909   test (ubuntu-latest)   completed        success      2021-09-04 12:41:10   2021-09-04 12:41:53   0.7166666666666667
3443854949   test (ubuntu-latest)   completed        success      2021-08-27 14:00:47   2021-08-27 14:01:36   0.8166666666666667
3428886296   test (ubuntu-latest)   completed        success      2021-08-26 03:51:14   2021-08-26 03:51:56   0.7
3422345982   test (ubuntu-latest)   completed        success      2021-08-25 13:20:52   2021-08-25 13:21:54   1.0333333333333334
3422049046   test (ubuntu-latest)   completed        success      2021-08-25 12:49:51   2021-08-25 12:50:33   0.7
3421981918   test (ubuntu-latest)   completed        success      2021-08-25 12:42:13   2021-08-25 12:42:59   0.7666666666666667
3421861740   test (ubuntu-latest)   completed        success      2021-08-25 12:28:41   2021-08-25 12:29:35   0.9
3421712060   test (ubuntu-latest)   completed        success      2021-08-25 12:12:20   2021-08-25 12:13:03   0.7166666666666667
3311170395   test (ubuntu-latest)   completed        success      2021-08-12 11:45:28   2021-08-12 11:46:17   0.8166666666666667
...

$ prstats actions jobs analyze  --owner itsubaki --repo q --format csv | column -t -s, | less -S
name                   start            end          runs_per_day          failure_rate          duration_avg(minutes)   duration_var(minutes)
test (ubuntu-latest)   2021-07-04       2021-07-11   0.14285714285714285   0                     0.9                     0
test (ubuntu-latest)   2021-07-11       2021-07-18   0                     0                     0                       0
test (ubuntu-latest)   2021-07-18       2021-07-25   0.14285714285714285   0                     0.8666666666666667      0
test (ubuntu-latest)   2021-07-25       2021-08-01   5.428571428571429     0.05263157894736842   0.7596491228070176      0.00839460724757415
test (ubuntu-latest)   2021-08-01       2021-08-08   0.2857142857142857    0                     0.825                   0.003402777777777773
test (ubuntu-latest)   2021-08-08       2021-08-15   0.5714285714285714    0                     0.7416666666666667      0.0024305555555555547
test (ubuntu-latest)   2021-08-15       2021-08-22   0                     0                     0                       0
test (ubuntu-latest)   2021-08-22       2021-08-29   1                     0                     0.8047619047619047      0.013231292517006813
test (ubuntu-latest)   2021-08-29       2021-09-05   0.14285714285714285   0                     0.7166666666666667      0
...
```