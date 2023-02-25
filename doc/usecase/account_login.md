# Login

## Actors

- Coach or Swimmer

## Stakeholders and Interests

- System: wants to know, who is interacting with it
- User: wants to interact with the system

## Description

User wants to login into system.

## Preconditions

1. User has created an account
2. User has verified the account

## Success Guarantees

User can interact with system.

## Success Scenario

1. User enters email
2. User enters password
3. User is logged in
4. User can interact with system

## Extensions

- 1A: User entered unknown email
  1.  Display generic message about invalid credentials
- 2A: User entered invalid password
  1.  Display generic message about invalid credentials
- 4A: User didn't verify account
  1.  Display message about not being verified
  2.  Don't allow user to interact with system
