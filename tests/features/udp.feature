Feature: send and receive UDP
  In order to ensure UDP connectivity
  Cyn should both send and receive UDP packets

  Scenario: one listen one send to nothing
    Given I run cyn --listen-udp 127.0.0.1:14563
    And I run cyn --send-udp 127.0.0.1:14568 --send-interval 10ms
    When I wait a moment
    Then the stdout contains "connection refused"

  Scenario: one listen one send (shorthand flags)
    Given I run cyn -u 127.0.0.1:14563
    And I run cyn -U 127.0.0.1:14563 --send-interval 10ms
    When I wait a moment
    Then the stdout contains "hi"

  Scenario: an instance is set to call itself via config file
    Given a configuration file that contains:
      """
      listen-udp:
        - 127.0.0.1:14568
      send-udp:
        - 127.0.0.1:14568
      send-interval: 10ms
      """
    And cyn is started with the config file
    When I wait 1 second
    Then the stdout contains "hi"
