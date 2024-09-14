import { NewTraining } from '~/api/generated';

const NEW_TRAINING_LOCAL_STORAGE_KEY = 'new-training';

export function saveTrainingToLocalStorage(training: NewTraining) {
  localStorage.setItem(
    NEW_TRAINING_LOCAL_STORAGE_KEY,
    JSON.stringify(training)
  );
}

export function loadTrainingFromLocalStorage() {
  const item = localStorage.getItem(NEW_TRAINING_LOCAL_STORAGE_KEY);
  if (!item) {
    return undefined;
  }
  const training = JSON.parse(item) as NewTraining;
  training.start = new Date(training.start);
  return training;
}

export function clearTrainingFromLocalStorage() {
  localStorage.removeItem(NEW_TRAINING_LOCAL_STORAGE_KEY);
}
