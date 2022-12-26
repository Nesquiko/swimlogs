# Create training session

## Actors

- Coach

## Description

Coach wants to create a training session.

## Success Guarantees

Session was saved and coach can see the traning session in collection of sessions.
Also coach can assign sessions to trainings.

## Success Scenario

1. Coach starts the session creation process
2. Coach inputs the day the session takes place
3. Coach inputs the start time of the session
4. Coach inputs the duration in minutes
5. Coach inputs a name for the session
6. Coach submits the session

## Extensions

- 5A: A session with specified name already exists
  1.  Highlight the name field
- 5B: User didn't input session name
  1.  A name in format `Day Start Duration` (ex. `Mon 17.00 60m`) is used as session name
