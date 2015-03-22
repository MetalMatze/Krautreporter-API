Feature: Test the articles endpoints

  Scenario: GET /articles
    When I request "GET" "/articles"
    Then I get a "200" response
    And the property "data" contains "268" items

  Scenario: GET /articles/17
    When I request "GET" "/articles/53"
    Then I get a "200" response
    And I scope into the property "data"
      And the property "id" is a integer equalling "53"
      And the property "title" is a string equalling "Israels Sicherheitsberater: Weshalb es nie Frieden geben wird"
      And the property "headline" is a string
      And the property "date" is a string equalling "2014-10-14"
      And the property "morgenpost" is a boolean equalling "false"
      And the property "url" is a string equalling "/53--der-himmel-bewahre-uns-davor-dass-wir-euch-europaer-ernst-nehmen"
      And the property "image" exists
      And the property "excerpt" is a string
      And the property "content" is a string
      And the property "author" is a integer equalling "13"
      And the property "created_at" is a string equalling "2015-03-21T12:31:49+0100"
      And the property "updated_at" is a string equalling "2015-03-21T13:17:09+0100"
