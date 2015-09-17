Feature: Test the authors endpoints

  Scenario: GET /authors
    When I request "GET" "/authors"
    Then I get a "200" response
    And the property "data" contains "10" items
    And I scope into the first property "data"
    And the property "id" is a integer equalling "1"
    And the property "name" is a string
    And the property "title" is a string
    And the property "url" is a string
    And the property "biography" is a string
    And the property "socialmedia" is a string
    And the property "created_at" is a iso8601 date
    And the property "updated_at" is a iso8601 date

  Scenario: GET /authors/123
    When I request "GET" "/authors/123"
    Then I get a "404" error response

  Scenario: GET /authors/5
    When I request "GET" "/authors/5"
    Then I get a "200" response
    And I scope into the property "data"
    And the property "id" is a integer equalling "5"
    And the property "name" is a string
    And the property "title" is a string
    And the property "url" is a string
    And the property "biography" is a string
    And the property "socialmedia" is a string
    And the property "created_at" is a iso8601 date
    And the property "updated_at" is a iso8601 date
