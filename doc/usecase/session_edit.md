# Edit training session

## Actors

- Coach

## Description

Coach wants to edit a training session.

## Preconditions

1. The session was created.

## Success Guarantees

Session was edited and coach can see it with new changes.

## Success Scenario

1. Coach starts editing process
2. Coach edits the session
3. Coach submits the edited session

## Extensions

- 3A: A session with specified name already exists
  1.  Highlight name field
- 3B: User edited name to empty
  1.  A name in format `Day Start Duration` (ex. `Mon 17.00 60m`) is used as session name
