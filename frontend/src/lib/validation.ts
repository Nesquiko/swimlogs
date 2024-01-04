import {
  InvalidTraining,
  InvalidTrainingSet,
  NewTraining,
  NewTrainingSet,
  StartType,
} from '../generated'

export const isInvalidTrainingEmpty = (it: InvalidTraining): boolean => {
  const isTrainingDataValid =
    it.durationMin === undefined && it.sets === undefined

  if (!isTrainingDataValid) return false

  let isValid = true
  it.invalidSets?.forEach((is) => {
    if (!isInvalidSetEmpty(is)) {
      isValid = false
    }

    if (is.subSets !== undefined) {
      is.subSets.forEach((sis) => {
        if (!isInvalidSetEmpty(sis)) {
          isValid = false
        }
      })
    }
  })
  return isValid
}

export const isInvalidSetEmpty = (is: InvalidTrainingSet | undefined) => {
  return (
    is?.setOrder === undefined &&
    is?.subSetOrder === undefined &&
    is?.repeat === undefined &&
    is?.distanceMeters === undefined &&
    is?.startType === undefined &&
    is?.startSeconds === undefined &&
    is?.subSets === undefined
  )
}

export const validateTraining = (t: NewTraining): InvalidTraining => {
  const invalidTraining: InvalidTraining = {}

  if (t.durationMin === undefined)
    invalidTraining.durationMin = "Duration can't be empty"
  else if (t.durationMin < 0)
    invalidTraining.durationMin = 'Duration must be positive'

  if (t.sets.length === 0)
    invalidTraining.sets = 'Training must have at least one set'

  const invalidSets = validateSets(t.sets)
  if (invalidSets) invalidTraining.invalidSets = invalidSets

  return invalidTraining
}

export const validateSets = (
  sets: Array<NewTrainingSet>
): Array<InvalidTrainingSet> | undefined => {
  const invalidSets: Array<InvalidTrainingSet> = []

  sets.forEach((set, index) => {
    invalidSets[index] = {}
    if (set.repeat <= 0) invalidSets[index].repeat = 'Repeat must be positive'

    if (set.distanceMeters !== undefined && set.distanceMeters <= 0)
      invalidSets[index].distanceMeters = 'Distance must be positive'

    if (!isStartTypeValid(set)) {
      invalidSets[index].startSeconds = 'Seconds must be positive'
    }

    if (set.subSets !== undefined) {
      invalidSets[index].subSets = validateSets(set.subSets)
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
