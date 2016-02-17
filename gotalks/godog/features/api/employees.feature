Feature: list yippee employees
  In order to manage employees
  As a product owner
  I need to be able to get employee list

  Background:
    Given there are employees:
      | John | Doe |
      | Jane | Doe |

  Scenario: I should be able to get employee list
    When I send "GET" request to "/users"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "users": [
          {"firstname": "John", "lastname": "Doe"},
          {"firstname": "Jane", "lastname": "Doe"}
        ]
      }
      """

  Scenario: should not allow POST method
    When I send "POST" request to "/users"
    Then the response code should be 405

  Scenario: should require authentication for protected resources
    When I send "GET" request to "/protected"
    Then the response code should be 401

  Scenario: should require authentication for protected resources
    When I send "GET" request to "/protected" as "user"
    Then the response code should be 200
