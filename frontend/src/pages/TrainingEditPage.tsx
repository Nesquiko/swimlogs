import { useTransContext } from '@mbarzda/solid-i18next';
import {
  createAsync,
  RouteSectionProps,
  useNavigate,
  useParams,
} from '@solidjs/router';
import {
  Component,
  createSignal,
  Match,
  onCleanup,
  onMount,
  Show,
  Switch,
} from 'solid-js';
import { createStore, SetStoreFunction, unwrap } from 'solid-js/store';
import { NewTrainingSet, ResponseError, Training } from 'swimlogs-api';
import { showToast } from '../App';
import { ToastMode } from '../components/DismissibleToast';
import { clearHeaderButton, setHeaderButton } from '../components/Header';
import TrainingEditButtonGroup from '../components/TrainingEditButtonGroup';
import TrainingPreview, {
  SkeletonTrainingPreview,
} from '../components/TrainingPreview';
import { cloneToSet } from '../lib/clone';
import {
  deleteSetInTraining,
  moveSetDownInTraining,
  moveSetUpInTraining,
  recalculateTotalDistance,
} from '../lib/training';
import { updateTrainingById } from '../state/trainings';
import EditSetPage from './EditSetPage';
import EditTrainingSessionPage from './EditTrainingSessionPage';
import { useOnBackcontext } from './Routing';

interface TrainingEditPageProps {
  trainingPromise: Promise<Training>;
}

const TrainingEditPage: Component<RouteSectionProps<TrainingEditPageProps>> = (
  props
) => {
  const params = useParams();
  const navigate = useNavigate();
  const [t] = useTransContext();
  const [onBack, setOnBack] = useOnBackcontext();

  const [showCreateSet, setShowCreateSet] = createSignal(false);
  const [editedSetIdx, setEditedSetIdx] = createSignal(-1);
  const [showTrainingSession, setShowTrainingSession] = createSignal(false);

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

  function onBackOverride() {
    if (showCreateSet()) {
      setShowCreateSet(false);
    } else if (editedSetIdx() !== -1) {
      setEditedSetIdx(-1);
    } else if (showTrainingSession()) {
      setShowTrainingSession(false);
    } else {
      setOnBack();
      onBack();
    }
  }

  onMount(() => {
    setOnBack(onBackOverride);

    setHeaderButton({
      icon: (
        <i
          classList={{
            'fa-calendar-days': !showTrainingSession(),
            'fa-person-swimming': showTrainingSession(),
          }}
          class="fa-solid fa-xl text-white cursor-pointer"
        />
      ),
      onClick: () => setShowTrainingSession(!showTrainingSession()),
    });
  });
  onCleanup(() => {
    setOnBack();
    clearHeaderButton();
  });

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

    const onCreateSet = (set: NewTrainingSet) => {
      set.setOrder = training.sets.length;
      const newSet = cloneToSet(set);
      newSet.id = crypto.randomUUID();
      setTraining('sets', [...training.sets, newSet]);
      setTraining('totalDistance', recalculateTotalDistance(training));
      setShowCreateSet(false);
    };

    const onEditSet = (set: NewTrainingSet) => {
      setTraining('sets', (sets) => {
        const tmp = sets[editedSetIdx()];
        set.setOrder = tmp.setOrder;
        const edit = cloneToSet(set);
        edit.id = tmp.id;
        sets[editedSetIdx()] = edit;
        return sets;
      });
      setTraining('totalDistance', recalculateTotalDistance(training));
      setEditedSetIdx(-1);
    };

    const onTrainingSessionSubmit = (trainingSession: {
      start: Date;
      durationMin: number;
    }) => {
      setTraining('start', trainingSession.start);
      setTraining('durationMin', trainingSession.durationMin);
      setShowTrainingSession(false);
    };

    return (
      <Switch fallback={trainingEditPreview(training, setTraining)}>
        <Match when={showCreateSet()}>
          <EditSetPage
            onSubmitSet={onCreateSet}
            submitLabel={t('add')}
            onCancel={() => setShowCreateSet(false)}
          />
        </Match>
        <Match when={editedSetIdx() !== -1}>
          <EditSetPage
            onSubmitSet={onEditSet}
            submitLabel={t('edit')}
            set={training.sets[editedSetIdx()]}
            onCancel={() => setEditedSetIdx(-1)}
          />
        </Match>
        <Match when={showTrainingSession()}>
          <EditTrainingSessionPage
            initial={{
              start: new Date(training.start),
              durationMin: training.durationMin,
            }}
            submitLabel={t('set.session')}
            onSubmit={onTrainingSessionSubmit}
            onBack={() => setShowTrainingSession(false)}
          />
        </Match>
      </Switch>
    );
  };

  const trainingEditPreview = (
    training: Training,
    setTraining: SetStoreFunction<Training>
  ) => {
    return (
      <>
        <TrainingPreview
          training={training}
          showSession={true}
          setOptions={[
            {
              text: t('edit'),
              icon: 'fa-pen',
              onClick: (setIdx) => setEditedSetIdx(setIdx),
            },
            {
              text: t('duplicate'),
              icon: 'fa-copy',
              onClick: (setIdx) => {
                const newSet = structuredClone(unwrap(training.sets[setIdx]));
                newSet.setOrder = training.sets.length;
                newSet.id = crypto.randomUUID();
                setTraining('sets', [...training.sets, newSet]);
                setTraining(
                  'totalDistance',
                  recalculateTotalDistance(training)
                );
              },
            },
            {
              text: t('move.up'),
              icon: 'fa-arrow-up',
              onClick: (setIdx) => moveSetUpInTraining(setIdx, setTraining),
              disabledFunc: (setIdx) => setIdx === 0,
            },
            {
              text: t('move.down'),
              icon: 'fa-arrow-down',
              onClick: (setIdx) =>
                moveSetDownInTraining(
                  setIdx,
                  training.sets.length,
                  setTraining
                ),
              disabledFunc: (setIdx) => setIdx === training.sets.length - 1,
            },
            {
              text: t('delete'),
              icon: 'fa-trash',
              onClick: (setIdx) => deleteSetInTraining(setIdx, setTraining),
            },
          ]}
        />

        <TrainingEditButtonGroup
          training={training}
          backLabel={t('back')}
          onBack={onBackOverride}
          onAddSet={() => setShowCreateSet(true)}
          confirmLabel={t('save')}
          onConfirm={() => updateTraining(training)}
        />
      </>
    );
  };

  return (
    <Show when={training()} fallback={<SkeletonTrainingPreview />}>
      {(t) => trainingEdit(t())}
    </Show>
  );
};

export default TrainingEditPage;
