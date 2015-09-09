Feature: Test the articles endpoints

  Scenario: GET /articles
    When I request "GET" "/articles"
    Then I get a "200" response
    And the property "data" contains "20" items

  Scenario: GET /articles/17
    When I request "GET" "/articles/17"
    Then I get a "200" response
