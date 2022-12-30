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

### Set

The most atomic unit in training. Each of them contains:

1. repeat multiplier
2. distance
3. what
4. starting rule (interval, pause, none)

#### Starting rule

Strategy to how to start in between the set repetitions.

Interval
: represented in form `I=XX`, where XX is how many seconds one repetition will take
Pause
: represented in form `P=XX`, where XX is how many seconds to wait in between the repetitions.
None
: just swim without any time restrictions

#### JSON representation

Simple set block with one repetition, with none as starting rule.

```json
{ "repeat": 1, "distance": 400, "what": "freestyle", "startingRule": "-" }
```

Set block with set of 8x50, with starting rule pause.

```json
{
  "repeat": 8,
  "distance": 50,
  "what": "BF/SF with fins and board",
  "startingRule": "P=20\""
}
```

Set block with set of 5x100, with starting rule interval.

```json
{
  "repeat": 5,
  "distance": 100,
  "what": "Mono or fins",
  "startingRule": "I=1'30\""
}
```

### Block

The whole training consists of blocks. Each of them have:

1. repeat multiplier
2. name
3. an array of [sets](#set), ordered by what to swim first

#### JSON representation

```json
{
  "repeat": 2,
  "name": "Main set",
  "sets": [
    {
      "repeat": 8,
      "distance": 100,
      "what": "fins, board, BF",
      "startingRule": "-"
    },
    {
      "repeat": 1,
      "distance": 100,
      "what": "cool down with board",
      "startingRule": "-"
    }
  ]
}
```

### Training

One training contains:

1. many [blocks](#block)
2. name
3. the date (in format `YYYY-MM-DD`) when it will occur.
4. sessionId ([see Session](./session.md)) to copy day, start time and duration to training
   - or day, start time and duration set manually

```json
{
  "date": "2022-04-04",
  "day": "tuesday",
  "startTime": "10:00",
  "durationMin": 120,
  "blocks": [
    {
      "repeat": 1,
      "name": "Warm up",
      "sets": [
        {
          "repeat": 1,
          "distance": 400,
          "what": "freestyle",
          "startingRule": "-"
        },
        {
          "repeat": 1,
          "distance": 400,
          "what": "fins, board",
          "startingRule": "-"
        }
      ]
    },
    {
      "repeat": 2,
      "name": "Main set",
      "sets": [
        {
          "repeat": 8,
          "distance": 100,
          "what": "fins, board, BF",
          "startingRule": "-"
        },
        {
          "repeat": 1,
          "distance": 100,
          "what": "cool down with board",
          "startingRule": "-"
        }
      ]
    }
  ]
}
```
