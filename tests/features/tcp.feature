Feature: send and receive TCP
  In order to ensure TCP connectivity
  Cyn should both send and receive on TCP connections

  Scenario: one tries to connect to nothing
    Given I run cyn --send-tcp 127.0.0.1:1234
    When I wait a moment
    Then the stdout contains "connection refused"

  Scenario: one listen one send (shorthand flags)
    Given I run cyn -t 127.0.0.1:15235
    And I run cyn -T 127.0.0.1:15235 -i 10ms
    When I wait a moment
    Then the stdout contains "connected"
    And the stdout contains "hi"

  Scenario: the connection is closed
    Given I run cyn -t 127.0.0.1:15236
    And I run cyn -T 127.0.0.1:15236 -i 10ms
    When I wait a moment
    And I stop process #1
    And I wait a moment
    Then the stdout contains "broken pipe"

  Scenario: the connection is reset
    Given I run cyn -t 127.0.0.1:15236
    And I run cyn -T 127.0.0.1:15236 -i 10ms
    When I wait a moment
    And I stop process #1
    And I wait a moment
    And I run cyn -t 127.0.0.1:15236
    And I wait a moment
    And I reset the output
    And I wait 1 second
    Then the stdout does not contain "broken pipe"

  Scenario: an instance is set to call itself via config file
    Given a configuration file that contains:
      """
      listen-tcp:
        - 127.0.0.1:24568
      send-tcp:
        - 127.0.0.1:24568
      send-interval: 10ms
      """
    When cyn is started with the config file
    And I wait 1 second
    Then the stdout contains "hi"

  Scenario: an instance is set to call itself via config file and the send interval is set via env variable
    Given a configuration file that contains:
      """
      listen-tcp:
        - 127.0.0.1:24568
      send-tcp:
        - 127.0.0.1:24568
      """
    And the environment variable CYNOMYS_SEND_INTERVAL is set to "10ms"
    When cyn is started with the config file
    And I wait a moment
    And I reset the output
    And I wait a moment
    Then the stdout contains "hi"
