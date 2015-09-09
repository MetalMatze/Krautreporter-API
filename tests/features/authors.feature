Feature: Test the authors endpoints

  Scenario: GET /authors
    When I request "GET" "/authors"
    Then I get a "200" response
    And the property "data" contains "10" items

  Scenario: GET /authors/5
    When I request "GET" "/authors/5"
    Then I get a "200" response
