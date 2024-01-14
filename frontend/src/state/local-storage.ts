import { NewTraining } from '../generated'

const LOCAL_STORAGE_KEY = 'new-training'

export const saveTrainingToLocalStorage = (training: NewTraining) => {
  localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(training))
}

export const loadTrainingFromLocalStorage = () => {
  const training = localStorage.getItem(LOCAL_STORAGE_KEY)
  if (!training) {
    return undefined
  }
  return JSON.parse(training) as NewTraining
}

export const clearTrainingFromLocalStorage = () => {
  localStorage.removeItem(LOCAL_STORAGE_KEY)
}
