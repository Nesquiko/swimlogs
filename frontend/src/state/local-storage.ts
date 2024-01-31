import { NewTraining } from 'swimlogs-api'

const NEW_TRAINING_LOCAL_STORAGE_KEY = 'new-training'

export const saveTrainingToLocalStorage = (training: NewTraining) => {
  localStorage.setItem(NEW_TRAINING_LOCAL_STORAGE_KEY, JSON.stringify(training))
}

export const loadTrainingFromLocalStorage = () => {
  const training = localStorage.getItem(NEW_TRAINING_LOCAL_STORAGE_KEY)
  if (!training) {
    return undefined
  }
  return JSON.parse(training) as NewTraining
}

export const clearTrainingFromLocalStorage = () => {
  localStorage.removeItem(NEW_TRAINING_LOCAL_STORAGE_KEY)
}
