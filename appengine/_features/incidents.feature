Feature:
    In order to get indicators of "Time to Restore Services" and "Change Failure Rate"
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

    Scenario: should create incident
        Given I set "Content-Type" header with "application/json"
        Given I set request body:
            """
            {
                "owner": "itsubaki",
                "repository": "ghz",
                "description": "[TEST] create incident",
                "sha": "abc",
                "resolved_at": "2021-12-05 10:00:00 UTC"
            }
            """
        When I send "POST" request to "/incidents/itsubaki/ghz"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "owner": "itsubaki",
                "repository": "ghz",
                "description": "[TEST] create incident",
                "sha": "abc",
                "resolved_at": "2021-12-05 10:00:00 UTC"
            }
            """

    Scenario: should fetch events
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/events"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/events"
            }
            """

    Scenario: should get TTR via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT owner, repository, pushed_at, resolved_at, TTR FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pushed`
            WHERE sha = "6f5dc2fc9b933ef6fd5f075924a5fec114405a25"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | pushed_at               | resolved_at             | TTR |
            | itsubaki | ghz        | 2021-12-24 09:01:29 UTC | 2021-12-24 10:01:29 UTC | 60  |

    Scenario: should get MTTR via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            WITH A AS(
            SELECT
            owner,
            repository,
            DATE(pushed_at) as date,
            PERCENTILE_CONT(TTR, 0.5) OVER(partition by DATE(pushed_at)) as MTTR
            FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pushed`
            )
            SELECT
            owner,
            repository,
            date,
            MAX(MTTR) as MTTR
            FROM A
            WHERE date = "2021-12-24"
            GROUP BY owner, repository, date
            """
        Then I get the following result:
            | owner    | repository | date       | MTTR |
            | itsubaki | ghz        | 2021-12-24 | 60   |

    Scenario: should get failure rate via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            WITH A AS (
            SELECT
            owner,
            repository,
            DATE(pushed_at) as date,
            COUNT(*) as failure
            FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pushed`
            GROUP BY date, owner, repository
            ), B AS (
            SELECT
            DATE(created_at) as date,
            COUNT(*) as pushed
            FROM `$PROJECT_ID.itsubaki_ghz.events_push`
            GROUP BY date
            )
            SELECT
            A.owner,
            A.repository,
            A.date,
            B.pushed,
            A.failure,
            A.failure / B.pushed as failure_rate
            FROM A
            INNER JOIN B
            ON A.date = B.date
            WHERE A.date = "2021-12-24"
            """
        Then I get the following result:
            | owner    | repository | date       | pushed | failure | failure_rate       |
            | itsubaki | ghz        | 2021-12-24 | 3      | 1       | 0.3333333333333333 |

    Scenario: should fetch pullreqs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/pullreqs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/pullreqs"
            }
            """

    Scenario: should fetch pullreqs/commits
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/pullreqs/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/pullreqs/commits"
            }
            """

    Scenario: should get TTR via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            SELECT owner, repository, merged_at, resolved_at, TTR FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs`
            WHERE sha = "aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | merged_at               | resolved_at             | TTR |
            | itsubaki | ghz        | 2021-12-08 09:41:12 UTC | 2021-12-08 10:41:12 UTC | 60  |

    Scenario: should get MTTR via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            WITH A AS(
            SELECT
            owner,
            repository,
            DATE(merged_at) as date,
            PERCENTILE_CONT(TTR, 0.5) OVER(partition by DATE(merged_at)) as MTTR
            FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs`
            )
            SELECT
            owner,
            repository,
            date,
            MAX(MTTR) as MTTR
            FROM A
            WHERE date = "2021-12-08"
            GROUP BY owner, repository, date
            """
        Then I get the following result:
            | owner    | repository | date       | MTTR |
            | itsubaki | ghz        | 2021-12-08 | 60   |

    Scenario: should get failure rate via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            WITH A AS (
            SELECT
            owner,
            repository,
            DATE(merged_at) as date,
            COUNT(*) as failure
            FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs`
            GROUP BY date, owner, repository
            ), B AS (
            SELECT
            DATE(merged_at) as date,
            COUNT(*) as merged
            FROM `$PROJECT_ID.itsubaki_ghz.pullreqs`
            WHERE state = "closed" AND merged_at != "0001-01-01 00:00:00 UTC"
            GROUP BY date
            )
            SELECT
            A.owner,
            A.repository,
            A.date,
            B.merged,
            A.failure,
            A.failure / B.merged as failure_rate
            FROM A
            INNER JOIN B
            ON A.date = B.date
            WHERE A.date = "2021-12-08"
            """
        Then I get the following result:
            | owner    | repository | date       | merged | failure | failure_rate |
            | itsubaki | ghz        | 2021-12-08 | 1      | 1       | 1.0          |
