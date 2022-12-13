Feature: run with nothing
  In order to use cyn more easily
  Cyn should provide help text when run with nothing given

  Scenario: run with no flags or config
    Given cyn is run with no flags or config
    When I wait a moment
    Then the stdout contains "no listeners or senders specified"
    And the stderr contains "Usage"
