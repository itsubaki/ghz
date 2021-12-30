Feature:
    In order to get indicators of incident
    As a DevOps practitioner

    Scenario: should get failure_rate and MTTR via pullrequest
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via PullRequest | aa0d19452f820c2088cbbe63d2fe2e18b67d3e4d | 2021-12-08 10:41:12 UTC |
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_pullreqs` WHERE date = "2021-12-08" LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | date       | merged | failure | failure_rate | MTTR |
            | itsubaki | ghz        | 2021-12-08 | 1      | 1       | 1.0          | 60.0 |

    Scenario: should get failure_rate and MTTR via commit
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | ghz        | [TEST] Incident via Commit | 6f5dc2fc9b933ef6fd5f075924a5fec114405a25 | 2021-12-24 10:01:29 UTC |
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._incidents_via_commits` WHERE date = "2021-12-24" LIMIT 1
            """
        Then I get the following result:
            | owner    | repository | date       | commits | failure | failure_rate       | MTTR |
            | itsubaki | ghz        | 2021-12-24 | 3       | 1       | 0.3333333333333333 | 60.0 |
