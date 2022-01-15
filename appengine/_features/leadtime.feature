Feature:
    In order to get indicators of "Lead Time for Changes"
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

    Scenario: should fetch commits
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/commits"
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

    Scenario: should get lead time via pushed
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._pushed_leadtime`
            WHERE sha = "25fd40317d3df7cafb770c3319fb122068724f25"
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | login    | message           | head_sha                                 | sha                                      | committed_at            | completed_at            | lead_time |
            | itsubaki | ghz        | 16163576    | tests         | itsubaki | Update some files | 4bb5472af191eff241ef8befdf88c44bed46ad85 | 25fd40317d3df7cafb770c3319fb122068724f25 | 2021-12-14 10:03:35 UTC | 2021-12-15 12:11:16 UTC | 1567      |

    Scenario: should get the median amount of lead time via pushed
        When I execute query with:
            """
            SELECT owner, repository, workflow_name, date, lead_time
            FROM `$PROJECT_ID.itsubaki_ghz._pushed_leadtime_median`
            WHERE date = "2021-12-30"
            """
        Then I get the following result:
            | owner    | repository | workflow_name | date       | lead_time |
            | itsubaki | ghz        | tests         | 2021-12-30 | 1.0       |

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

    Scenario: should get lead time via pullrequests
        When I execute query with:
            """
            SELECT * FROM `$PROJECT_ID.itsubaki_ghz._pullreqs_leadtime`
            WHERE sha = "d80f4a0921f36da81b2d27a8d27d4328ada988c8"
            """
        Then I get the following result:
            | owner    | repository | workflow_id | workflow_name | pullreq_id | pullreq_number | login    | title                     | message                   | merge_commit_sha                         | sha                                      | committed_at            | completed_at            | lead_time |
            | itsubaki | ghz        | 16163576    | tests         | 811825160  | 6              | itsubaki | Update google api version | Update google api version | 2566237cf6179830721e4357eb53089d53598fa7 | d80f4a0921f36da81b2d27a8d27d4328ada988c8 | 2021-12-30 07:34:38 UTC | 2021-12-30 07:38:53 UTC | 4         |

    Scenario: should get the median amount of lead time via pullrequests
        When I execute query with:
            """
            SELECT owner, repository, workflow_name, date, lead_time
            FROM `$PROJECT_ID.itsubaki_ghz._pullreqs_leadtime_median`
            WHERE date = "2021-12-30"
            """
        Then I get the following result:
            | owner    | repository | workflow_name | date       | lead_time |
            | itsubaki | ghz        | tests         | 2021-12-30 | 4.0       |

