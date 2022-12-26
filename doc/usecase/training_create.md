# Create training

## Actors

- Coach

## Description

Coach wants to create a training.

## Success Guarantees

Training was saved and coach can see it in some collection.

## Success Scenario

1. Coach starts the creation process
2. Coach defines the training
3. Coach gives training a name
4. Coach assigns training session to the training
   - Or coach assigns date, time and duration to the training
5. Coach submits the training

## Extensions

- 5A System isn't responding
  1.  Stay in the creation process
- 5B Training with given name exists
  1.  Highlight the name field
- 5C Training isn't in valid format
  1.  Highlight invalid fields
