# ghz

[![PkgGoDev](https://pkg.go.dev/badge/github.com/itsubaki/ghz)](https://pkg.go.dev/github.com/itsubaki/ghz)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsubaki/ghz?style=flat-square)](https://goreportcard.com/report/github.com/itsubaki/ghz)


## Installation

```shell
go install github.com/itsubaki/ghz@latest
```

## Examples

```shell
ghz commits fetch --owner itsubaki --repository ghz
ghz commits list  --owner itsubaki --repository ghz --format csv | column -t -s, | less -S
```

```tsv
sha                                        login      date                  message
20838b4c89bfc3b018c82cd5f290dd014a92cc15   itsubaki   2022-06-30 13:32:59   Update description                                                     
51d552207c68303387b8628043f08b4a2eb62743   itsubaki   2022-06-27 03:59:37   Update README.md                                                       
23ea5269430a5849d6efc08f294001594499ba47   itsubaki   2022-06-27 03:58:29   Merge branch 'main' of https://github.com/itsubaki/ghz                 
e7297e67c2ff550c1f7dcbdfb99b3111ba514692   itsubaki   2022-06-27 03:58:04   Remove pkg dir                                                         
1217efbaf31cbc201341c139f41a3207edab5dd5   itsubaki   2022-06-11 06:16:56   Update README.md                                                               
```

```shell
ghz pullreqs fetch --owner itsubaki --repository ghz
ghz pullreqs list  --owner itsubaki --repository ghz --format csv | column -t -s, | less -S
```

```tsv
id          number   title                       login      state    created_at            updated_at            merged_at             closed_at             merge_commit_sha               
888469760   9        Update Makefile             itsubaki   closed   2022-03-24 14:47:03   2022-03-24 14:47:16   2022-03-24 14:47:12   2022-03-24 14:47:12   9b79b1fb4e2ad0cff8873cfc70f7077
831019272   7        Update cron schedule        itsubaki   closed   2022-01-25 03:01:27   2022-02-08 04:20:27   2022-01-25 03:30:33   2022-01-25 03:30:33   166b2e6654700e6fc900cc9323e0e0b
811825160   6        Update google api version   itsubaki   closed   2021-12-30 07:35:04   2021-12-30 07:37:28   2021-12-30 07:37:25   2021-12-30 07:37:25   2566237cf6179830721e4357eb53089
799869157   3        Update cron.yaml            itsubaki   closed   2021-12-10 12:20:33   2021-12-16 03:17:52   null                  2021-12-10 12:55:19   1c8b23ea3d1f49fb0f05b7080a27d65
797653180   2        Update some files           itsubaki   closed   2021-12-08 09:41:04   2021-12-08 09:41:15   2021-12-08 09:41:12   2021-12-08 09:41:12   7d849f71d16ab268c5cf51f97c6147d
788450216   1        Update some files           itsubaki   closed   2021-11-25 02:29:53   2021-11-25 02:30:01   2021-11-25 02:30:00   2021-11-25 02:30:00   5ac595c8277ad5898972c31dfdbe710
```

```shell
ghz pullreqs commits fetch --owner itsubaki --repository ghz
ghz pullreqs commits list  --owner itsubaki --repository ghz --format csv | column -t -s, | less -S
```

```tsv
id          number   sha                                        login      date                  message                     
888469760   9        95f21a055aaf62130d1a6f4d29b0f4a0e92d638b   itsubaki   2022-03-24 14:46:42   Update Makefile             
831019272   7        765c9d33721fa282df4e03879c499c0a8961ac89   itsubaki   2022-01-25 03:00:54   Update cron schedule        
811825160   6        d80f4a0921f36da81b2d27a8d27d4328ada988c8   itsubaki   2021-12-30 07:34:38   Update google api version   
811825160   6        40b702fed5fb318558c9f83e7445b9fe25446220   itsubaki   2021-12-30 07:37:07   Fix typo                    
799869157   3        cb40954bca00a65f3a8ee8b85a827052e4edbb29   itsubaki   2021-12-10 12:20:15   Update cron.yaml            
797653180   2        2b2889d032ef3801428e5b51333d02f6e14de88c   itsubaki   2021-12-08 09:39:39   Update some files           
797653180   2        762ac8925298e6f13f345c7fab3192e1124fa96c   itsubaki   2021-12-08 01:47:56   Update pullreqs             
797653180   2        aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d   itsubaki   2021-12-08 04:19:52   Update some files           
797653180   2        55a6923b72ce762fe23363313e833e8278ddeeba   itsubaki   2021-12-08 05:17:59   Update pullreqs             
797653180   2        8c82a7d6637bf284b58cdbe00605bf00d327e9e8   itsubaki   2021-12-08 06:50:43   Update some files           
797653180   2        a3c493cf5010ab604cdc05d528672c54f0b47167   itsubaki   2021-12-08 09:03:25   Update some files           
788450216   1        6c53bf63b37ff24328e02508be6b86885a9b1e68   itsubaki   2021-11-25 02:29:30   Update some files  
```

```shell
ghz actions runs fetch --owner itsubaki --repository ghz
ghz actions runs list  --owner itsubaki --repository ghz --format csv | column -t -s, | less -S
```

```tsv
workflow_id   workflow_name   run_id       run_number   status      conclusion   created_at            updated_at            head_commit.sha                            head_commit.date    
16163576      tests           2590287280   214          completed   success      2022-06-30 13:33:20   2022-06-30 13:33:45   20838b4c89bfc3b018c82cd5f290dd014a92cc15   2022-06-30 13:32:59 
16163576      tests           2566678851   213          completed   success      2022-06-27 03:59:45   2022-06-27 04:00:13   51d552207c68303387b8628043f08b4a2eb62743   2022-06-27 03:59:37 
16163576      tests           2479017958   212          completed   success      2022-06-11 06:16:57   2022-06-11 06:17:27   1217efbaf31cbc201341c139f41a3207edab5dd5   2022-06-11 06:16:56 
16163576      tests           2478766501   211          completed   success      2022-06-11 04:32:58   2022-06-11 04:33:17   a2c2919022929eacbe13ca19559e5429e439b3b0   2022-06-11 04:32:57 
16163576      tests           2439942555   210          completed   success      2022-06-04 14:39:49   2022-06-04 14:40:08   463b997b845a5d388b18f9466f1bce9c6d60ebea   2022-06-04 14:39:36 
16163576      tests           2428981101   209          completed   success      2022-06-02 14:44:58   2022-06-02 14:45:26   63dde2f6359f422332ba36c23ae037e902306973   2022-06-02 14:43:21 
16163576      tests           2428973336   208          completed   success      2022-06-02 14:43:38   2022-06-02 14:44:13   63dde2f6359f422332ba36c23ae037e902306973   2022-06-02 14:43:21 
16163576      tests           2425894270   207          completed   success      2022-06-02 03:44:06   2022-06-02 03:53:23   d72d0968269e805a0d9c7ab0af129e3b7014cc5c   2022-06-02 03:44:05 
16163576      tests           2425877361   206          completed   success      2022-06-02 03:38:39   2022-06-02 03:47:34   f50f3b3d884524dab2fb276ee3bb76e203586fac   2022-06-02 03:38:28 
```

```shell
ghz actions jobs fetch --owner itsubaki --repository ghz
ghz actions jobs list  --owner itsubaki --repository ghz --format csv | column -t -s, | less -S
```

```tsv
run_id       job_id       job_name               status      conclusion   started_at            completed_at          
2590287280   7132461391   test (ubuntu-latest)   completed   success      2022-06-30 13:33:32   2022-06-30 13:33:43   
2566678851   7065991966   test (ubuntu-latest)   completed   success      2022-06-27 03:59:53   2022-06-27 04:00:11   
2479017958   6841469129   test (ubuntu-latest)   completed   success      2022-06-11 06:17:05   2022-06-11 06:17:26   
2478766501   6840982976   test (ubuntu-latest)   completed   success      2022-06-11 04:33:06   2022-06-11 04:33:16   
2439942555   6738856399   test (ubuntu-latest)   completed   success      2022-06-04 14:39:56   2022-06-04 14:40:06   
2428981101   6710544303   test (ubuntu-latest)   completed   success      2022-06-02 14:45:07   2022-06-02 14:45:24   
2428973336   6710521156   test (ubuntu-latest)   completed   success      2022-06-02 14:43:45   2022-06-02 14:44:11   
2425894270   6701950535   test (ubuntu-latest)   completed   success      2022-06-02 03:47:41   2022-06-02 03:53:21   
2425877361   6701901811   test (ubuntu-latest)   completed   success      2022-06-02 03:41:36   2022-06-02 03:47:32   
2425865939   6701847393   test (ubuntu-latest)   completed   success      2022-06-02 03:35:10   2022-06-02 03:41:26   
2419099054   6682423516   test (ubuntu-latest)   completed   success      2022-06-01 03:35:42   2022-06-01 03:41:47   
```

```shell
ghz pullreqs list --owner itsubaki --repository ghz | jq -r '[.id, .number, .user.login, .title, .state] | @csv' | tr -d '"' | column -t -s, | less -S
```

```tsv
888469760  9  itsubaki  Update Makefile            closed
831019272  7  itsubaki  Update cron schedule       closed
811825160  6  itsubaki  Update google api version  closed
799869157  3  itsubaki  Update cron.yaml           closed
797653180  2  itsubaki  Update some files          closed
788450216  1  itsubaki  Update some files          closed
```
