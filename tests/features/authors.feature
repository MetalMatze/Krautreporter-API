Feature: Test the authors endpoints

  Scenario: GET /authors
    When I request "GET" "/authors"
    Then I get a "200" response
    And the property "data" contains "67" items

  Scenario: GET /authors/13
    When I request "GET" "/authors/13"
    Then I get a "200" response
    And I scope into the property "data"
      And the property "id" is a integer equalling "13"
      And the property "name" is a string equalling "Tilo Jung"
      And the property "title" is a string equalling "Politik"
      And the property "url" is a string equalling "/13--tilo-jung"
      And the property "image" is a string equalling "/system/user/profile_image/13/thumb_retina_KUcC_EJQiDrfoGvsdsSz0SjfNhTEouNaONYGGC1vaMw.png"
      And the property "biography" is a string equalling "Tilo Jung"
      And the property "socialmedia" is a string
      And the property "created_at" is a string equalling "2015-03-21T12:31:42+0100"
      And the property "updated_at" is a string equalling "2015-03-21T12:35:22+0100"
