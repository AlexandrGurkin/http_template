Feature: get version
  In order to know godog version
  As an API user
  I need to be able to request version

  Scenario: does not allow POST method
    When I send request to version handler
    Then the response code should be 200
    And the response should contain field: "version"