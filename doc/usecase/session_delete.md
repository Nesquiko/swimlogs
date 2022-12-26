# Delete training session

## Actors

- Coach

## Description

Coach wants to delete training session.

## Success Guarantees

Session was deleted and it is no longer in collection of sessions.

## Success Scenario

1. Coach deletes the session
2. Coach confirms the deletion

## Extensions

- 2A: Session doesn't exits
  1.  Refresh the session collection
