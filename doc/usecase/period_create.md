# Create training period

## Actors

- Coach

## Description

Coach wants to create a training period, which starts and ends at some specific
dates.

## Success Guarantees

Period was saved and coach can see new period in some collection of periods.

## Success Scenario

1. Coach starts the period creation process.
2. Coach inputs the start date.
3. Coach inputs the end date.
4. Coach inputs the name of the period.
5. Coach submits the period.

## Extensions

- 5A: System isn't responding
  1.  Inform user
  2.  Stay in the creation process
- 5B: Period with given name exists
  1.  Don't submit the period
  2.  Highlight the name field
- 5C: End date is before start date
  1.  Don't submit the training
  2.  Highlight end date part
