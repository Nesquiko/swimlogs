import { NewTrainingSet } from 'swimlogs-api'

export function cloneSet(s: NewTrainingSet) {
  const newSet = {} as NewTrainingSet
  newSet.repeat = s.repeat
  newSet.setOrder = s.setOrder
  newSet.startType = s.startType
  newSet.description = s.description
  newSet.startSeconds = s.startSeconds
  newSet.distanceMeters = s.distanceMeters
  newSet.totalDistance = s.totalDistance

  newSet.equipment = []
  s.equipment?.forEach((e) => {
    newSet.equipment?.push(e)
  })
  return newSet
}
