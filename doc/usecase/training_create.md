# Create training

## Actors

- Coach

## Description

Coach wants to create a training.

## Success Guarantees

Training was saved and coach can see the training in some collection.

## Success Scenario

1. Coach starts the creation process.
2. Coach defines the training.
3. Coach gives training a name.
4. Coach assigns training session to the training
5. Coach submits the training.

## Extensions

- 5A System isn't responding
  1.  Inform user
  2.  Stay in the creation process
- 5B Training with given name exists
  1.  Don't submit the training
  2.  Highlight the name field
- 5C Training isn't in valid format
  1.  Don't submit the training
  2.  Highlight invalid part
