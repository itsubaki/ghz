# Using Google AppEngine, BigQuery

- Automatically fetch repository metadata and insert into BigQuery
- Provides some views

## Required

- Google AppEngine
- Google BigQuery

## Configuration

- `cron.yaml`
- `secrets.yaml`

```shell
$ cat secrets.yaml
env_variables:
  PAT: YOUR_GITHUB_PERSONAL_ACCESS_TOKEN
```

## Deploy

```shell
$ gcloud beta app deploy app.yaml cron.yaml
```
