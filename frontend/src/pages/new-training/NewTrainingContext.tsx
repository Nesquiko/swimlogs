import { createContext, ParentComponent, useContext } from 'solid-js';
import { createStore, SetStoreFunction } from 'solid-js/store';
import { NewTraining } from '~/api/generated';
import { nowWithHoursInRange } from '~/lib/datetime';

export const StartTimeHours = [
  '06',
  '07',
  '08',
  '09',
  '10',
  '11',
  '12',
  '13',
  '14',
  '15',
  '16',
  '17',
  '18',
  '19',
  '20',
];

export const StartTimeMinutes = ['00', '15', '30', '45'];

const NEW_TRAINING_LOCAL_STORAGE_KEY = 'new-training';

function saveTrainingToLocalStorage(training: NewTraining) {
  localStorage.setItem(
    NEW_TRAINING_LOCAL_STORAGE_KEY,
    JSON.stringify(training)
  );
}

function loadTrainingFromLocalStorage() {
  const item = localStorage.getItem(NEW_TRAINING_LOCAL_STORAGE_KEY);
  if (!item) {
    return undefined;
  }
  const training = JSON.parse(item) as NewTraining;
  training.start = new Date(training.start);
  return training;
}

export function makeNewTrainingContext() {
  const [training, _setTraining] = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: nowWithHoursInRange(
        parseInt(StartTimeHours[0]),
        parseInt(StartTimeHours[StartTimeHours.length - 1])
      ),
      durationMinutes: 60,
      sets: [],
    }
  );

  // @ts-expect-error: Not assignable…
  const setTraining: SetStoreFunction<NewTraining> = (...params: [any]) => {
    _setTraining(...params);
    saveTrainingToLocalStorage(training);
  };

  return [training, setTraining] as const;
}

type NewTrainingContextType = ReturnType<typeof makeNewTrainingContext>;
const NewTrainingContext = createContext<NewTrainingContextType>();
export const useNewTraining = () => useContext(NewTrainingContext);

export const NewTrainingContextProvider: ParentComponent = (props) => {
  return (
    <NewTrainingContext.Provider value={makeNewTrainingContext()}>
      {props.children}
    </NewTrainingContext.Provider>
  );
};
