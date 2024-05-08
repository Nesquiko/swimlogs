import { useTransContext } from '@mbarzda/solid-i18next';
import { useNavigate } from '@solidjs/router';
import { createStore } from 'solid-js/store';
import { showToast } from '../App';
import { ToastMode } from '../components/DismissibleToast';
import { NewTraining } from 'swimlogs-api';
import {
  clearTrainingFromLocalStorage,
  loadTrainingFromLocalStorage,
  saveTrainingToLocalStorage,
} from '../state/local-storage';
import { addTrainingDetail, trainingApi } from '../state/trainings';
import { Component } from 'solid-js';
import { StartTimeHours } from '../components/SessionEditForm';
import { nowWithHoursInRange } from '../lib/datetime';
import TrainingEditProcess from '../components/TrainingEditProcess';

const TrainingCreatePage: Component = () => {
  const [t] = useTransContext();
  const navigate = useNavigate();

  const [training, setTraining] = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: nowWithHoursInRange(
        parseInt(StartTimeHours[0]),
        parseInt(StartTimeHours[StartTimeHours.length - 1])
      ),
      durationMin: 60,
      totalDistance: 0,
      sets: [],
    }
  );

  const onSubmit = () => {
    trainingApi
      .createTraining({ newTraining: training })
      .then((res) => {
        addTrainingDetail(res);
        showToast(t('training.created', 'Training created'));
        clearTrainingFromLocalStorage();
      })
      .catch((e) => {
        console.error('error', e);
        showToast(t('training.creation.error'), ToastMode.ERROR);
      })
      .finally(() => {
        navigate('/', { replace: true });
      });
  };

  return (
    <div class="px-4">
      <TrainingEditProcess
        training={training}
        setTraining={setTraining}
        onSubmit={onSubmit}
        onUpdate={(tr) => saveTrainingToLocalStorage(tr)}
        onDelete={() => {
          clearTrainingFromLocalStorage();
          navigate('/', { replace: true });
        }}
      />
    </div>
  );
};

export default TrainingCreatePage;
