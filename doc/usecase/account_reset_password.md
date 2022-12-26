# Reest account password

## Actors

- Coach or Swimmer

## Description

User wants to reset password.

## Preconditions

1. User must have an account.

## Success Guarantees

User can login with new password.

## Success Scenario

1. User selects reset password option
2. User enters an email
3. User uses received reset password email to navigate to reset password process
4. User enters new password
5. User re-enters new password
6. User submits new password

## Extensions

- 2A: User entered an unknown email
  1.  Show that email was sent to entered address
- 3A: Email doesn't come
  1.  User can retry resetting password
- 5A: Passwords don't match
  1.  Highlight password fields
