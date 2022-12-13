Feature: send and receive UDP
  In order to ensure UDP connectivity
  Cyn should both send and receive UDP packets

  Scenario: one listen one send to nothing
    Given cyn is listening for UDP on 127.0.0.1:14563
    And cyn is sending UDP to 127.0.0.1:14568
    When I wait 2 seconds
    Then the stdout contains "connection refused"

  Scenario: one listen one send
    Given cyn is listening for UDP on 127.0.0.1:14563
    And cyn is sending UDP to 127.0.0.1:14563
    When I wait 2 seconds
    Then the stdout contains "hi"
