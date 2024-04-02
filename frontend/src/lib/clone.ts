import { NewTrainingSet, TrainingSet } from 'swimlogs-api'

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

export function cloneToSet(ns: NewTrainingSet): TrainingSet {
  const newSet = {} as TrainingSet

  newSet.id = ''
  newSet.repeat = ns.repeat
  newSet.setOrder = ns.setOrder
  newSet.startType = ns.startType
  newSet.description = ns.description
  newSet.startSeconds = ns.startSeconds
  newSet.distanceMeters = ns.distanceMeters
  newSet.totalDistance = ns.totalDistance
  newSet.group = ns.group

  newSet.equipment = []
  ns.equipment?.forEach((e) => {
    newSet.equipment?.push(e)
  })
  return newSet
}
