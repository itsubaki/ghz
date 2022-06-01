# Using Google AppEngine, BigQuery

- Automatically fetch repository metadata and insert into BigQuery
- Provides some views

## DORA's Four Keys

 - [Definitions](https://github.com/GoogleCloudPlatform/fourkeys/blob/main/METRICS.md)

### Deployment Frequency

 - The number of deployments to production per week.
 - Elite: Over the last 3 months, the median number is equal to or greater than 3.


### Lead Time for Changes
 
 - The time for a commit to be deployed into production.
 - Elite: Over the last 3 months, the median amount is equal to or less than 1 day.

### Time to Restore Services

 - The time between the <deployment> which caused the failure and the remediation.
   - In this product, `resolved_at - merged_at` or `resolved_at - pushed_at`, To be exact `resolved_at - deployed_at`
 - Elite: Over the last 3 months, the median amount is equal to or less than 1 day.

### Change Failure Rate

 - The number of failures per the number of deployments.
   - In this product, `count(failure)/count(merged)`, To be exact `count(failure)/count(depyloyed)`.
 - Elite: Over the last 3 months, the median amount is equal to or less than 15 %.

## Required

- Google AppEngine
- Google BigQuery
- Google Cloud Scheduler

## Configuration

- `cron.yaml`
- `secrets.yaml`

```shell
$ cat secrets.yaml
env_variables:
  PAT: YOUR_GITHUB_PERSONAL_ACCESS_TOKEN
```

## Deploying to AppEngine

```shell
$ gcloud app deploy app.yaml cron.yaml
```

## Integration Tests

- Create google cloud service account.
- Put `credentials.json` to root directory.

```shell
$ make test
GOOGLE_APPLICATION_CREDENTIALS=../credentials.json go test ./appengine -v -coverprofile=coverage.out -covermode=atomic --godog.format=pretty -coverpkg ./...

Feature:
  In order to fetch repository metadata
  As an X-Appengine-Cron
  I need to be able to fetch request

  Background:
    Given I set "X-Appengine-Cron" header with "true"         # features_test.go:65 -> *apiFeature

  Scenario: should fetch commits                              # _features/fetch.feature:9
    When I send "GET" request to "/_fetch/itsubaki/q/commits" # features_test.go:48 -> *apiFeature
    Then the response code should be 200                      # features_test.go:57 -> *apiFeature


1 scenarios (1 passed)
3 steps (3 passed)
4.211440948s
testing: warning: no tests to run
PASS
coverage: 15.1% of statements in ./...
ok      github.com/itsubaki/ghz/appengine       4.459s  coverage: 15.1% of statements in ./... [no tests to run]
```

## GitHub Actions

- Set repository secrets
  - `GOOGLE_APPLICATION_CREDENTIALS`: json of `github@${PROJECT_ID}.iam.gserviceaccount.com`
  - `PAT`: Your GitHub Personal Access Token

## IAM & Admin

- `${PROJECT_ID}@appspot.gserviceaccount.com` (App Engine default service account)
  - App Engine Admin
  - BigQuery Admin
  - Cloud Scheduler Job Runner
  - Cloud Trace Agent
  - Cloud Profiler Agent
- `github@${PROJECT_ID}.iam.gserviceaccount.com` (CI/CD @ GitHub Actions )
  - App Engine Admin
  - BigQuery Admin
  - Cloud Build Editor
  - Cloud Scheduler Admin
  - Storage Object Admin
  - Service Account User
- `localhost@${PROJECT_ID}.iam.gserviceaccount.com` (Integration Tests @ localhost)
  - BigQuery Admin
- `${PROFJECT_NUMBER}@cloudbuild.gserviceaccount.com`
  - Cloud Build Service Account
