Feature:
    In order to get indicators of "Lead Time for Changes"
    As a DevOps practitioner


    Scenario: should fetch actions runs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/runs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/runs",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch commits
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/commits",
                "next_token": "@string@"
            }
            """

    Scenario: should fetch events
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/events"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/events",
                "next_token": "@string@"
            }
            """

    Scenario: should get lead time via commits

    Scenario: should fetch pullreqs
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/pullreqs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/pullreqs",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch pullreqs/commits
        Given I set "X-Appengine-Cron" header with "true"
        When I send "GET" request to "/_fetch/itsubaki/ghz/pullreqs/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/pullreqs/commits",
                "next_token": "@number@"
            }
            """

    Scenario: should get lead time via pullrequests
