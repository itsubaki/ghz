Feature:
    In order to get indicators of "Time to Restore Services" and "Change Failure Rate"
    As a DevOps practitioner

    Scenario: should create dataset
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/_init?renew=true"
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

    Scenario: should get approximated TTR via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT owner, repository, pushed_at, resolved_at, TTR
            FROM `$PROJECT_ID.itsubaki_ghz._pushed_ttr`
            WHERE sha = "6f5dc2fc9b933ef6fd5f075924a5fec114405a25"
            """
        Then I get the following result:
            | owner    | repository | pushed_at               | resolved_at             | TTR |
            | itsubaki | ghz        | 2021-12-24 09:01:29 UTC | 2021-12-24 10:01:29 UTC | 60  |

    Scenario: should get approximated MTTR via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT owner, repository, date, MTTR
            FROM `$PROJECT_ID.itsubaki_ghz._pushed_ttr_median`
            WHERE date = "2021-12-24"
            """
        Then I get the following result:
            | owner    | repository | date       | MTTR |
            | itsubaki | ghz        | 2021-12-24 | 60   |

    Scenario: should get TTR via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            WITH A AS (
            SELECT head_sha, resolved_at
            FROM `$PROJECT_ID.itsubaki_ghz.incidents` as A
            INNER JOIN `$PROJECT_ID.itsubaki_ghz._pushed` as B
            ON A.sha = B.sha
            )
            SELECT B.owner, B.repository, B.workflow_name, B.updated_at as changed_at, A.resolved_at, TIMESTAMP_DIFF(A.resolved_at, B.updated_at, MINUTE) as TTR
            FROM A
            INNER JOIN `$PROJECT_ID.itsubaki_ghz.workflow_runs` as B
            ON A.head_sha = B.head_sha
            WHERE Date(resolved_at) = "2021-12-24"
            """
        Then I get the following result:
            | owner    | repository | workflow_name | changed_at              | resolved_at             | TTR |
            | itsubaki | ghz        | tests         | 2021-12-24 09:02:30 UTC | 2021-12-24 10:01:29 UTC | 58  |

    Scenario: should get failure rate via pushed
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Pushed | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT owner, repository, date, pushed, failure, failure_rate
            FROM `$PROJECT_ID.itsubaki_ghz._pushed_failure_rate`
            WHERE date = "2021-12-24"
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

    Scenario: should get approximated TTR via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            SELECT owner, repository, merged_at, resolved_at, TTR
            FROM `$PROJECT_ID.itsubaki_ghz._pullreqs_ttr`
            WHERE sha = "aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d"
            """
        Then I get the following result:
            | owner    | repository | merged_at               | resolved_at             | TTR |
            | itsubaki | ghz        | 2021-12-08 09:41:12 UTC | 2021-12-08 10:41:12 UTC | 60  |

    Scenario: should get approximated MTTR via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            SELECT owner, repository, date, MTTR
            FROM `$PROJECT_ID.itsubaki_ghz._pullreqs_ttr_median`
            WHERE date = "2021-12-08"
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
            SELECT owner, repository, date, merged, failure, failure_rate
            FROM `$PROJECT_ID.itsubaki_ghz._pullreqs_failure_rate`
            WHERE date = "2021-12-08"
            """
        Then I get the following result:
            | owner    | repository | date       | merged | failure | failure_rate |
            | itsubaki | ghz        | 2021-12-08 | 1      | 1       | 1.0          |
