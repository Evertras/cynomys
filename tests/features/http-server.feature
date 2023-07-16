Feature: host HTTP server for live status
  In order to easily see the status of my network
  Cyn should be able to run an HTTP server with the latest information
  So that I can visually check if data is flowing

  Scenario: http server is enabled
    Given I run cyn --listen-tcp 127.0.0.1:24184 --http.address 127.0.0.1:24185 --listen-udp 127.0.0.1:24183
    And I wait a moment
    # TODO: better verification, this is quick and dirty for now
    Then the page at http://127.0.0.1:24185 contains "table"
    And the page at http://127.0.0.1:24185 contains "<td>127.0.0.1:24184</td>"
    And the page at http://127.0.0.1:24185 contains "<td>127.0.0.1:24183</td>"
