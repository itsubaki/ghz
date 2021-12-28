Feature:
    In order to fetch repository metadata
    As an X-Appengine-Cron
    I need to be able to fetch request

    Background:
        Given I set "X-Appengine-Cron" header with "true"

    Scenario: should fetch commits
        When I send "GET" request to "/_fetch/itsubaki/q/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/commits",
                "next_token": "@string@"
            }
            """

    Scenario: should fetch pullreqs
        When I send "GET" request to "/_fetch/itsubaki/q/pullreqs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/pullreqs",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch pullreqs/commits
        When I send "GET" request to "/_fetch/itsubaki/q/pullreqs/commits"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/pullreqs/commits",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch events
        When I send "GET" request to "/_fetch/itsubaki/q/events"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/events",
                "next_token": "@string@"
            }
            """

    Scenario: should fetch releases
        When I send "GET" request to "/_fetch/itsubaki/q/releases"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/releases",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch incidents
        When I send "GET" request to "/_fetch/itsubaki/q/incidents"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/incidents"
            }
            """

    Scenario: should fetch actions runs
        When I send "GET" request to "/_fetch/itsubaki/q/actions/runs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/actions/runs",
                "next_token": "@number@"
            }
            """

    Scenario: should fetch actions jobs
        When I send "GET" request to "/_fetch/itsubaki/q/actions/jobs"
        Then the response code should be 200
        Then the response should match json:
            """
            {
                "path": "/_fetch/itsubaki/q/actions/jobs",
                "next_token": "@number@"
            }
            """
