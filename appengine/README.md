# Using Google AppEngine, BigQuery

- Automatically fetch repository metadata and insert into BigQuery
- Provides some views

## DORA's Four Keys

- Deployment Frequency: The number of deployments to production per day.
- Lead Time for Changes: The median amount of time for a commit to be deployed into production.
- Time to Restore Services: The median amount of time between the <deployment> which caused the failure and the remediation.
  - In this product, The median amount of time between the <push event/pull request merge> which caused the failure and the remediation.
- Change Failure Rate: The number of failures per the number of deployments.

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

## Deployment

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
- `github@${PROJECT_ID}.iam.gserviceaccount.com` (CI/CD @ GitHub Actions )
  - App Engine Admin
  - BigQuery Admin
  - Cloud Build Editor
  - Cloud Scheduler Admin
  - Service Account User
  - Storage Object Admin
- `localhost@${PROJECT_ID}.iam.gserviceaccount.com` (Integration Tests @ localhost)
  - BigQuery Admin
- `${PROFJECT_NUMBER}@cloudbuild.gserviceaccount.com`
  - Cloud Build Service Account
