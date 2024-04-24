import { NewTraining } from 'swimlogs-api';

const NEW_TRAINING_LOCAL_STORAGE_KEY = 'new-training';

export const saveTrainingToLocalStorage = (training: NewTraining) => {
  localStorage.setItem(
    NEW_TRAINING_LOCAL_STORAGE_KEY,
    JSON.stringify(training)
  );
};

export const loadTrainingFromLocalStorage = () => {
  const item = localStorage.getItem(NEW_TRAINING_LOCAL_STORAGE_KEY);
  if (!item) {
    return undefined;
  }
  const training = JSON.parse(item) as NewTraining;
  training.start = new Date(training.start);
  return training;
};

export const clearTrainingFromLocalStorage = () => {
  localStorage.removeItem(NEW_TRAINING_LOCAL_STORAGE_KEY);
};
