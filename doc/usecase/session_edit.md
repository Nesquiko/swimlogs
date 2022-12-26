# Edit training session

## Actors

- Coach

## Description

Coach wants to edit training session.

## Preconditions

1. The session must exits.

## Success Guarantees

Session was edited and coach can see it in some collection of session.

## Success Scenario

1. Coach starts editing process.
2. Coach edits the session.
3. Coach submits the edited session.

## Extensions

- 3A: System isn't responding
  1.  Inform user
  2.  Stay in the editing process
- 3B: An equal session already exists
  1.  Don't submit the session
- 3C: Session contains invalid data
  1.  Highlight invalid fields
