Feature:
    In order to get indicators of "Deployment Frequency"
    As a DevOps practitioner

    Scenario: should create dataset
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/_init"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/_init"
            }
            """

    Scenario: should fetch actions runs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/runs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/runs"
            }
            """

    Scenario: should fetch actions runs/update
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/runs/update"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/runs/update"
            }
            """

    Scenario: should get deployment frequency via runs
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._frequency_runs`
            WHERE date = "2021-12-25"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | date       | runs | duration_avg       |
            | itsubaki | ghz        | 16163576    | tests         | 2021-12-25 | 7    | 0.5714285714285714 |

    Scenario: should fetch actions jobs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/jobs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/jobs"
            }
            """

    Scenario: should fetch actions jobs/update
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/jobs/update"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/jobs/update"
            }
            """

    Scenario: should get deployment frequency via jobs
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._frequency_jobs`
            WHERE date = "2021-12-15"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | job_name             | date       | runs | duration_avg |
            | itsubaki | ghz        | 16163576    | tests         | test (ubuntu-latest) | 2021-12-15 | 4    | 0.25         |

    Scenario: should fetch releases
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/releases"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/releases"
            }
            """

    Scenario: should get deployment frequency via releases
        When I execute query with:
            """
            SELECT
            owner, repository, Date(published_at) as date, count(name) as releases
            FROM `$PROJECT_ID.itsubaki_ghz.releases`
            GROUP BY owner, repository, date
            ORDER BY date
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | date       | releases |
            | itsubaki | ghz        | 2021-12-30 | 1        |
