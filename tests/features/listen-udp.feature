Feature: listen for UDP
  In order to ensure UDP connectivity
  Cyn should acknowledge UDP packets

  Scenario: nothing is sent
    Given cyn is listening for UDP on 127.0.0.1:14563
    When I wait 1 second
    Then there is no output

  Scenario: a single UDP packet is sent
    Given cyn is listening for UDP on 127.0.0.1:14564
    When I send a UDP packet containing "hello" to 127.0.0.1:14564
    Then the stdout contains "hello"

  Scenario: multiple UDP packets are sent
    Given cyn is listening for UDP on 127.0.0.1:14564
    When I send a UDP packet containing "hello" to 127.0.0.1:14564
    And I send a UDP packet containing "another" to 127.0.0.1:14564
    Then the stdout contains "hello"
    And the stdout contains "another"
