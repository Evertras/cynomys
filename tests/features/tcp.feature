Feature: send and receive TCP
  In order to ensure TCP connectivity
  Cyn should both send and receive on TCP connections

  Scenario: one tries to connect to nothing
    Given cyn is sending TCP to 127.0.0.1:1234
    When I wait a moment
    Then the stdout contains "connection refused"

  Scenario: one listen one send
    Given cyn is listening for TCP on 127.0.0.1:15235
    And cyn is sending TCP to 127.0.0.1:15235
    When I wait 2 seconds
    Then the stdout contains "connected"
    And the stdout contains "hi"

  Scenario: an instance is set to call itself via config file
    Given a configuration file that contains:
      """
      listen-tcp:
        - 127.0.0.1:24568
      send-tcp:
        - 127.0.0.1:24568
      """
    And cyn is started with the config file
    When I wait 2 seconds
    Then the stdout contains "hi"
