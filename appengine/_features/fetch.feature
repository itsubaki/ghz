Feature:
    In order to fetch repository metadata
    As an X-Appengine-Cron
    I need to be able to fetch request

    Background:
        Given I set "X-Appengine-Cron" header with "true"

    Scenario: should fetch commits
        When I send "GET" request to "/_fetch/itsubaki/ghz/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/commits",
                "next_token": "@string@"
            }
            """

    Scenario: should fetch pullreqs
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
        When I send "GET" request to "/_fetch/itsubaki/ghz/pullreqs/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/pullreqs/commits",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch events
        When I send "GET" request to "/_fetch/itsubaki/ghz/events"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/events",
                "next_token": "@string@"
            }
            """

    Scenario: should fetch releases
        When I send "GET" request to "/_fetch/itsubaki/ghz/releases"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/releases",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch incidents
        When I send "GET" request to "/_fetch/itsubaki/ghz/incidents"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/incidents"
            }
            """

    Scenario: should fetch actions runs
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/runs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/runs",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch actions jobs
        When I send "GET" request to "/_fetch/itsubaki/ghz/actions/jobs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/ghz/actions/jobs",
                "next_token": "@number@"
            }
            """
