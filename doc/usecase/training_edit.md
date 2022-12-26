# Edit training

## Actors

- Coach

## Description

Coach wants to edit training.

## Preconditions

1. The training must exist.

## Success Guarantees

Training was edited and coach can see it in collection of trainings.

## Success Scenario

1. Coach starts editing process
2. Coach edits the training
3. Coach submits the edited training

## Extensions

- 3A: System isn't responding
  1.  Inform user
- 3B Training with given name exists
  1.  Highlight the name field
- 3C Training isn't in valid format
  1.  Highlight invalid fields
