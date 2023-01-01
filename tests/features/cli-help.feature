Feature: run to trigger CLI help
  In order to use cyn more easily
  Cyn should provide help text when run with incorrect configurations

  Scenario: run with no flags or config
    Given cyn is run with no flags or config
    When I wait a moment
    Then the stdout contains "no listeners or senders specified"
    And the stderr contains "Usage"

  Scenario: run with unknown flags
    Given cyn is run with an unknown flag
    When I wait a moment
    Then the stderr contains "Usage"
