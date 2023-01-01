Feature: listen for TCP
  In order to ensure TCP connectivity
  Cyn should accept incoming TCP connections

  Scenario: nothing is sent
    Given cyn is listening for TCP on 127.0.0.1:24563
    When I wait 1 second
    Then there is no output

  Scenario: a TCP connection is made with no data sent
    Given cyn is listening for TCP on 127.0.0.1:24564
    When I connect with TCP to 127.0.0.1:24564
    And I wait a moment
    Then the stdout contains "TCP connected"

  Scenario: a TCP connection is made and then disconnected
    Given cyn is listening for TCP on 127.0.0.1:24565
    When I connect with TCP to 127.0.0.1:24565
    And I wait a moment
    And I disconnect my TCP connection
    And I wait a moment
    Then the stdout contains "TCP connected"
    And the stdout contains "TCP disconnected"

  Scenario: a TCP connection is made and data is sent
    Given cyn is listening for TCP on 127.0.0.1:24565
    When I connect with TCP to 127.0.0.1:24565
    And I send "my tcp stuff" over my TCP connection
    And I wait a moment
    Then the stdout contains "TCP connected"
    And the stdout contains "my tcp stuff"
