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
5. Coach submits the session

## Extensions

- 5A: System isn't responding
  1.  Stay in the creation process
- 5B: An equal session already exists
  1.  Don't submit the session
- 5C: Session contains invalid data
  1.  Highlight invalid fields
