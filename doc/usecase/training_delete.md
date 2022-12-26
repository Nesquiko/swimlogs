# Delete training

## Actors

- Coach

## Description

Coach wants to delete training.

## Preconditions

1. The training must exist

## Success Guarantees

Training was deleted and it is no longer in collection of trainings.

## Success Scenario

1. Coach deletes the training
2. Coach confirms the deletion

## Extensions

- 2A: System isn't responding
  1.  Inform user
- 2B: Training doesn't exits
  1.  Refresh the training collection
