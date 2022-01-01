Feature:
    In order to get indicators of "Time to Restore Services" and "Change Failure Rate"
    As a DevOps practitioner

    Scenario: should create incidents table
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/incidents"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/incidents"
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

    Scenario: should get failure_rate and MTTR via commits
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Commit | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_commits`
            WHERE date = "2021-12-24"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | date       | pushed | failure | failure_rate | MTTR |
            | itsubaki | ghz        | 2021-12-24 | 2      | 1       | 0.5          | 60.0 |

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

    Scenario: should get failure_rate and MTTR via pullrequests
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs`
            WHERE date = "2021-12-08"
            LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | date       | merged | failure | failure_rate | MTTR |
            | itsubaki | ghz        | 2021-12-08 | 1      | 1       | 1.0          | 60.0 |

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
