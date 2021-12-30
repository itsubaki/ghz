Feature:
    In order to get indicators of incident
    As a BigQuery User
    I need to be able to incidents request

    Scenario: should get failure_rate and MTTR via pullrequest
        Given the following incidents exist:
            | owner    | repository | description                     | sha                                      | resolved_at             |
            | itsubaki | q          | [TEST] Incident via PullRequest | 7b2619e89065d96e683d70a72512e2883c1a2cf6 | 2021-07-30 13:04:37 UTC |
        When I execute query with:
            """
            "SELECT * FROM `$PROJECT_ID.itsubaki_q._incidents_via_pullreqs` WHERE date = \"2021-07-30\" LIMIT 1"
            """
        Then I get the following result:
            | owner    | repository | date       | merged | failure | failure_rate | MTTR |
            | itsubaki | q          | 2021-07-30 | 1      | 1       | 1.0          | 60.0 |

    Scenario: should get failure_rate and MTTR via commit
        Given the following incidents exist:
            | owner    | repository | description                | sha                                      | resolved_at             |
            | itsubaki | q          | [TEST] Incident via Commit | ad79208ce9ad1fce87b298ae28c6c518dc2a0486 | 2021-12-26 15:31:05 UTC |
        When I execute query with:
            """
            "SELECT * FROM `$PROJECT_ID.itsubaki_q._incidents_via_commits` WHERE date = \"2021-12-26\" LIMIT 1"
            """
        Then I get the following result:
            | owner    | repository | date       | commits | failure | failure_rate | MTTR |
            | itsubaki | q          | 2021-12-26 | 2       | 1       | 0.5          | 60.0 |
