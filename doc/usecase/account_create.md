# Create user account

## Actors

- Coach or Swimmer

## Stakeholders and Interests

- System: wants to know identities of users interacting with it
- User: wants to interact with the system

## Description

User wants to create an account in order to interact with the system. To create
an account user enters email and password.

## Success Guarantees

An account was created and user can login and begin verification process.

## Success Scenario

1. User starts account creation process
2. User enters email
3. User enters password
4. User re-enters password
5. User submits the credentials

## Extensions

- 2A: Email is already registered
  1.  Highlight email field
- 2B: Email is invalid
  1.  Highlight email field
- 4A: Passwords don't match
  1.  Highlight password fields
