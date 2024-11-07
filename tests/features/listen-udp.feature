Feature: listen for UDP
  In order to ensure UDP connectivity
  Cyn should acknowledge UDP packets

  Scenario: nothing is sent
    #Given cyn is listening for UDP on 127.0.0.1:14563
    Given I run cyn --listen.udp 127.0.0.1:14563
    When I wait 1 second
    Then there is no output

  Scenario: a single UDP packet is sent
    Given I run cyn --listen.udp 127.0.0.1:14564
    When I send a UDP packet containing "hello" to 127.0.0.1:14564
    Then the stdout contains "hello"

  Scenario: multiple UDP packets are sent (shorthand flag)
    Given I run cyn -u 127.0.0.1:14564
    When I send a UDP packet containing "hello" to 127.0.0.1:14564
    And I send a UDP packet containing "another" to 127.0.0.1:14564
    Then the stdout contains "hello"
    And the stdout contains "another"

  Scenario: the udp listener is set via config file
    Given a configuration file that contains:
      """
      listen:
        udp:
          - 127.0.0.1:14568
      """
    And cyn is started with the config file
    When I send a UDP packet containing "hello" to 127.0.0.1:14568
    Then the stdout contains "hello"
