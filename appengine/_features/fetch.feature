Feature:
    In order to fetch repository metadata
    As an X-Appengine-Cron
    I need to be able to fetch request

    Background:
        Given I set "X-Appengine-Cron" header with "true"

    Scenario: should fetch commits
        When I send "GET" request to "/_fetch/itsubaki/q/commits"
        Then the response code should be 200

    Scenario: should fetch pullreqs
        When I send "GET" request to "/_fetch/itsubaki/q/pullreqs"
        Then the response code should be 200
