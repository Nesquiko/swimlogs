# How does a swim training plan look like

## Current form of trainings

Current trainings look something like this:

| Distance | What                      | Start   |
| :------- | :------------------------ | :------ |
| 400      | freestyle                 | -       |
| 8x50     | BF/SF with fins and board | I=1'    |
| 200      | SF with Mono              | -       |
| 100      | cool down                 | -       |
| 15x100   | Fins or Mono              | I=1'30" |
| 100      | cool down                 | -       |
| 15x100   | Fins or Mono              | I=1'30" |
| 100      | cool down                 | -       |

## New form of trainings

New form would consist of N `set blocks`. Couple JSON examples of `set block`:

Simple set block with one repetition.

```json
{ "distance": "400", "what": "freestyle", "start": "-" }
```

Set block with set of 8x50.

```json
{
  "distance": "8x50",
  "what": "BF/SF with fins and board",
  "start": "P=20\""
}
```

Set block with set of 5x100

```json
{
  "distance": "5x100",
  "what": "Mono or fin",
  "start": "I=1'30\""
}
```

## JSON of one training

One training contains many blocks and the date (in format `YYYY-MM-DD`) when it will occur.

```json
{
  "date": "2022-04-04",
  "training": [
    { "distance": "400", "what": "freestyle", "start": "-" },
    {
      "distance": "8x50",
      "what": "BF/SF with fins and board",
      "start": "P=20\""
    },
    { "distance": "100", "what": "cool down", "start": "-" },
    {
      "distance": "5x100",
      "what": "Mono or fin",
      "start": "I=1'30\""
    },
    { "distance": "100", "what": "cool down", "start": "-" }
  ]
}
```

Also it contains either:

1. An id of training session from which to copy day, start time and duration

```json
{
  "sessionId": "b78e1672-90d3-44aa-8019-c802307948e4"
}
```

2. Or day, start time and duration

```json
{
  "day": "tuesday",
  "startTime": "10:00",
  "durationMin": 120
}
```
