import { useTransContext } from '@mbarzda/solid-i18next';
import {
  createAsync,
  RouteSectionProps,
  useNavigate,
  useParams,
} from '@solidjs/router';
import { Component, Show } from 'solid-js';
import { createStore } from 'solid-js/store';
import { ResponseError, Training } from 'swimlogs-api';
import { showToast } from '../App';
import { ToastMode } from '../components/DismissibleToast';
import { SkeletonSetsPreview } from '../components/SetsPreview';
import TrainingEditProcess from '../components/TrainingEditProcess';
import { updateTrainingById } from '../state/trainings';

interface TrainingEditPageProps {
  trainingPromise: Promise<Training>;
}

const TrainingEditPage: Component<RouteSectionProps<TrainingEditPageProps>> = (
  props
) => {
  const params = useParams();
  const navigate = useNavigate();
  const [t] = useTransContext();

  const training = createAsync(() =>
    props.data!.trainingPromise.catch((e: ResponseError) => {
      let msg = t('server.error');
      if (e.response?.status === 404) {
        msg = t('training.not.found');
      }
      showToast(msg, ToastMode.ERROR);
      navigate('/', { replace: true });
      return Promise.reject(e);
    })
  );

  async function updateTraining(training: Training) {
    updateTrainingById(params.id, training)
      .then(() => showToast(t('training.edited')))
      .catch((e: ResponseError) => {
        let msg = t('server.error');
        if (e.response?.status === 404) {
          msg = t('training.not.found');
        }
        showToast(msg, ToastMode.ERROR);
      })
      .finally(() => navigate('/', { replace: true }));
  }

  const trainingEdit = (tr: Training) => {
    const [training, setTraining] = createStore<Training>(structuredClone(tr));
    // TODO test this

    return (
      <TrainingEditProcess<Training>
        training={training}
        // @ts-ignore
        setTraining={setTraining}
        onSubmit={updateTraining}
        labels={{ setsPreviewLabel: 'edit.training.sets' }}
      />
    );
  };

  return (
    <div class="px-4">
      <Show when={training()} fallback={<SkeletonSetsPreview />}>
        {(t) => trainingEdit(t())}
      </Show>
    </div>
  );
};

export default TrainingEditPage;
