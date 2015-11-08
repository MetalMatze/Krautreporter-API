Feature: Test the articles endpoints

  Scenario: GET /articles
    When I request "GET" "/articles"
    Then I get a "200" response
    And the property "data" contains "20" items
    And I scope into the first property "data"
    And the property "id" is an integer equalling "101"
    And the property "order" is an integer equalling "100"
    And the property "title" is a string
    And the property "headline" is a string
    And the property "date" is a iso8601 date
    And the property "morgenpost" is a boolean equalling "false"
    And the property "preview" is a boolean
    And the property "url" is a string
    And the property "excerpt" is a string
    And the property "content" is a string
    And the property "author_id" is an integer
    And the property "created_at" is a iso8601 date
    And the property "updated_at" is a iso8601 date

  Scenario: GET /articles/123
    When I request "GET" "/articles/123"
    Then I get a "404" error response

  Scenario: GET /articles/17
    When I request "GET" "/articles/17"
    Then I get a "200" response
    And I scope into the property "data"
    And the property "id" is an integer equalling "17"
    And the property "order" is an integer equalling "16"
    And the property "title" is a string
    And the property "headline" is a string
    And the property "date" is a iso8601 date
    And the property "morgenpost" is a boolean equalling "false"
    And the property "preview" is a boolean
    And the property "url" is a string
    And the property "excerpt" is a string
    And the property "content" is a string
    And the property "author_id" is an integer
    And the property "created_at" is a iso8601 date
    And the property "updated_at" is a iso8601 date
