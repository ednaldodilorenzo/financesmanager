const { Given, When, Then } = require("cypress-cucumber-preprocessor/steps");

Given('I open the login page', () => {
  cy.visit('/login'); // Adjust the URL to your application's login page
});

When('I submit login credentials', () => {
  cy.get('input[name="username"]').type('user'); // Replace with your form's username and password fields
  cy.get('input[name="password"]').type('password');
  cy.get('form').submit();
});

Then('I should see the homepage', () => {
  cy.url().should('include', '/home'); // Adjust according to your application's home page URL
});
