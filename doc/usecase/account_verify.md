# Verify user account

## Actors

- Coach or Swimmer

## Stakeholders and Interests

- System: wants to verify the identity of the user
- User: wants to interact with the system

## Description

After creating an user account, a verification email is sent to email address
specified by the user during creation of the account.

## Preconditions

1. User created an account

## Success Guarantees

User can interact with system after login.

## Success Scenario

1. User receives a verification email
2. User clicks verify
3. User is verified and can interact with system

## Extensions

- 1A: User didn't receive the verification email
  1.  New verification email is sent
- 2A: Verification failed
  1.  Inform user
  2.  Keep the account in unverified state
  3.  User can't interact with the system, even if login was successful
