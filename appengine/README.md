# Using Google AppEngine, BigQuery

- Automatically fetch repository metadata and insert into BigQuery
- Provides some views

## Required

- Google AppEngine
- Google BigQuery
- Google Cloud Scheduler

## Configuration and Deploy

- `cron.yaml`
- create `secrets.yaml`

```shell
$ cat secrets.yaml
env_variables:
  PAT: YOUR_GITHUB_PERSONAL_ACCESS_TOKEN
```

```shell
$ gcloud beta app deploy app.yaml cron.yaml
```

## Integration Test

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

- Set repository secrets `GOOGLE_APPLICATION_CREDENTIALS` and `PAT`
