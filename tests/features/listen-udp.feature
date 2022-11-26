Feature: listen for UDP
  In order to ensure UDP connectivity
  Cyn should acknowledge UDP packets

  Scenario: nothing is sent
    Given cyn is listening for UDP on 127.0.0.1:14563
    When I wait 1 second
    Then there is no output
