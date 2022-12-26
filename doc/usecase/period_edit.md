# Edit training period

## Actors

- Coach

## Description

Coach wants to edit training period.

## Preconditions

1. The period must exits.

## Success Guarantees

Period was edited and coach can see it in some collection of periods.

## Success Scenario

1. Coach starts editing process.
2. Coach edits the period.
3. Coach submits the edited period.

## Extensions

- 3A: System isn't responding
  1.  Inform user
  2.  Stay in the editing process
- 3B: Period with given name exists
  1.  Don't submit the edit
  2.  Highlight the name field
- 3C: End date is before start date
  1.  Don't submit the period
  2.  Highlight end date part
