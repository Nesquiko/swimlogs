import {
  InvalidTraining,
  InvalidTrainingSet,
  NewTraining,
  NewTrainingSet,
  StartType
} from '../generated'

export const validateTraining = (
  t: NewTraining
): InvalidTraining | undefined => {
  const invalidTraining: InvalidTraining = {}

  if (t.durationMin === undefined)
    invalidTraining.durationMin = "Duration can't be empty"
  else if (t.durationMin < 0)
    invalidTraining.durationMin = 'Duration must be positive'

  const invalidSets = validateSets(t.sets)
  if (invalidSets) invalidTraining.invalidSets = invalidSets

  if (Object.keys(invalidTraining).length === 0) return undefined

  return invalidTraining
}

function validateSets(
  sets: Array<NewTrainingSet>
): Array<InvalidTrainingSet> | undefined {
  const invalidSets: Array<InvalidTrainingSet> = []

  sets.forEach((set, index) => {
    invalidSets[index] = {}
    if (set.repeat <= 0) invalidSets[index].repeat = 'Repeat must be positive'

    if (set.distanceMeters !== undefined && set.distanceMeters <= 0)
      invalidSets[index].distanceMeters = 'Distance must be positive'

    if (!isStartTypeValid(set)) {
      invalidSets[index].startSeconds = 'Seconds must be positive'
    }
  })

  return invalidSets
}

function isStartTypeValid(set: NewTrainingSet): boolean {
  if (set.startType === StartType.None) {
    return true
  }

  // else it is Interval or Pause
  if (set.startSeconds === undefined) {
    return false
  }

  return set.startSeconds > 0
}
