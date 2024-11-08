Feature: watch latency in stdout
  In order to check latency in an ad hoc fashion
  Cyn should send the latency to stdout

  Scenario: no latency in stdout by default
    Given I run cyn -u 127.0.0.1:14567
    And I run cyn -U 127.0.0.1:14567 --send.interval 200ms
    When I wait 1 second
    Then the stdout does not contain "latency"

  Scenario: latency is displayed in stdout when asked for
    Given I run cyn -u 127.0.0.1:14567
    And I run cyn -U 127.0.0.1:14567 --send.interval 200ms --sinks.stdout.enabled
    When I wait 1 second
    Then the stdout contains "latency"

  Scenario: latency is displayed in stdout when configured via file
    Given a configuration file that contains:
      """
      send:
        interval: 10ms
        udp:
          - 127.0.0.1:14567
      sinks:
        stdout:
          enabled: true
      """
    And I run cyn -u 127.0.0.1:14567
    When cyn is started with the config file 
    When I wait 1 second
    Then the stdout contains "latency"
