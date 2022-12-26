# Delete training session

## Actors

- Coach

## Description

Coach wants to delete training session.

## Preconditions

1. The session must exist

## Success Guarantees

Session was deleted and it is no longer in collection of session.

## Success Scenario

1. Coach deletes the session
2. Coach confirms the deletion

## Extensions

- 2A: System isn't responding
  1.  Inform user
- 2B: Session doesn't exits
  1.  Refresh the session collection
