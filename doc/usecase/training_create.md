# Create training

## Actors

- Coach

## Description

Coach wants to create a training.

## Success Guarantees

Training was saved and coach can see it in collection.

## Success Scenario

1. Coach starts the creation process
2. Coach defines the training
3. Coach gives training a name
4. Coach assigns training session to the training
5. Coach submits the training

## Extensions

- 3A: Training with given name exists
  1.  Highlight the name field
- 4A: User doesn't assign a training session
  1.  User is prompted to assigns date, time and duration to the training
- 5B: Training isn't in valid format
  1.  Highlight invalid fields
