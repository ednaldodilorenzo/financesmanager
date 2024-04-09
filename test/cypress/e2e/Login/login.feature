Feature: Login functionality

Scenario: Successful login
  Given I open the login page
  When I submit login credentials
  Then I should see the homepage
