import { produce, SetStoreFunction } from 'solid-js/store'
import { NewTraining, Training } from 'swimlogs-api'

export function recalculateTotalDistance(
  training: NewTraining | Training
): number {
  let newTotalDistance = 0
  training.sets.forEach(
    (s) => (newTotalDistance += s.repeat * (s.distanceMeters ?? 0))
  )
  return newTotalDistance
}

export function moveSetUpInNewTraining(
  setIdx: number,
  setTraining: SetStoreFunction<NewTraining>
): void {
  if (setIdx === 0) return

  setTraining('sets', (sets) => {
    const tmp = sets[setIdx]
    sets[setIdx] = sets[setIdx - 1]
    sets[setIdx - 1] = tmp
    return sets.map((s, i) => ({ ...s, setOrder: i }))
  })
}

export function moveSetUpInTraining(
  setIdx: number,
  setTraining: SetStoreFunction<Training>
): void {
  if (setIdx === 0) return

  setTraining('sets', (sets) => {
    const tmp = sets[setIdx]
    sets[setIdx] = sets[setIdx - 1]
    sets[setIdx - 1] = tmp
    return sets.map((s, i) => ({ ...s, setOrder: i }))
  })
}

export function moveSetDownInNewTraining(
  setIdx: number,
  setCount: number,
  setTraining: SetStoreFunction<NewTraining>
) {
  if (setIdx === setCount - 1) return

  setTraining('sets', (sets) => {
    const tmp = sets[setIdx]
    sets[setIdx] = sets[setIdx + 1]
    sets[setIdx + 1] = tmp
    return sets.map((s, i) => ({ ...s, setOrder: i }))
  })
}

export function moveSetDownInTraining(
  setIdx: number,
  setCount: number,
  setTraining: SetStoreFunction<Training>
) {
  if (setIdx === setCount - 1) return

  setTraining('sets', (sets) => {
    const tmp = sets[setIdx]
    sets[setIdx] = sets[setIdx + 1]
    sets[setIdx + 1] = tmp
    return sets.map((s, i) => ({ ...s, setOrder: i }))
  })
}

export function deleteSetInNewTraining(
  setIdx: number,
  setTraining: SetStoreFunction<NewTraining>
) {
  setTraining(
    produce((t) => {
      t.totalDistance = t.totalDistance - t.sets[setIdx].totalDistance
      return t
    })
  )
  setTraining('sets', (sets) =>
    sets.filter((_, i) => i !== setIdx).map((s, i) => ({ ...s, setOrder: i }))
  )
}

export function deleteSetInTraining(
  setIdx: number,
  setTraining: SetStoreFunction<Training>
) {
  setTraining(
    produce((t) => {
      t.totalDistance = t.totalDistance - t.sets[setIdx].totalDistance
      return t
    })
  )
  setTraining('sets', (sets) =>
    sets.filter((_, i) => i !== setIdx).map((s, i) => ({ ...s, setOrder: i }))
  )
}
