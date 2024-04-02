import { useTransContext } from '@mbarzda/solid-i18next';
import { useNavigate } from '@solidjs/router';
import { createStore } from 'solid-js/store';
import { showToast } from '../App';
import { ToastMode } from '../components/DismissibleToast';
import { NewTraining, NewTrainingSet } from 'swimlogs-api';
import {
  clearTrainingFromLocalStorage,
  loadTrainingFromLocalStorage,
  saveTrainingToLocalStorage,
} from '../state/local-storage';
import { addTrainingDetail, trainingApi } from '../state/trainings';
import {
  Component,
  createSignal,
  Match,
  onCleanup,
  onMount,
  Switch,
} from 'solid-js';
import EditSetPage from './EditSetPage';
import {
  deleteSetInNewTraining,
  moveSetDownInNewTraining,
  moveSetUpInNewTraining,
  recalculateTotalDistance,
} from '../lib/training';
import TrainingPreview from '../components/TrainingPreview';
import { cloneSet } from '../lib/clone';
import EditTrainingSessionPage from './EditTrainingSessionPage';
import DismissibleModal from '../components/DismissibleModal';
import TrainingEditButtonGroup from '../components/TrainingEditButtonGroup';
import { useOnBackcontext } from './Routing';

const NewTrainingPage: Component = () => {
  const [t] = useTransContext();
  const navigate = useNavigate();
  const [onBack, setOnBack] = useOnBackcontext();

  const [showTrainingSession, setShowTrainingSession] = createSignal(false);
  const [showCreateSet, setShowCreateSet] = createSignal(false);
  const [editedSetIdx, setEditedSetIdx] = createSignal(-1);

  const [openConfirmationModal, setOpenConfirmationModal] = createSignal(false);
  const [training, setTraining] = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: new Date(new Date().setHours(18, 0, 0, 0)),
      durationMin: 60,
      totalDistance: 0,
      sets: [],
    }
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

  onMount(() => setOnBack(onBackOverride));
  onCleanup(() => setOnBack());

  const onSubmit = (training: NewTraining) => {
    trainingApi
      .createTraining({ newTraining: training })
      .then((res) => {
        addTrainingDetail(res);
        showToast(t('training.created', 'Training created'));
      })
      .catch((e) => {
        console.error('error', e);
        showToast(t('training.creation.error'), ToastMode.ERROR);
      })
      .finally(() => {
        navigate('/', { replace: true });
      });
  };

  const onTrainingSessionSubmit = (trainingSession: {
    start: Date;
    durationMin: number;
  }) => {
    setShowTrainingSession(false);
    setTraining('start', trainingSession.start);
    setTraining('durationMin', trainingSession.durationMin);
    onSubmit(training);
    clearTrainingFromLocalStorage();
  };

  const onCreateSet = (set: NewTrainingSet) => {
    setShowCreateSet(false);
    set.setOrder = training.sets.length;
    setTraining('sets', [...training.sets, set]);
    setTraining('totalDistance', recalculateTotalDistance(training));
    saveTrainingToLocalStorage(training);
  };

  const onEditSet = (set: NewTrainingSet) => {
    setTraining('sets', (sets) => {
      const tmp = sets[editedSetIdx()];
      set.setOrder = tmp.setOrder;
      sets[editedSetIdx()] = set;
      return sets;
    });
    setEditedSetIdx(-1);
    setTraining('totalDistance', recalculateTotalDistance(training));
    saveTrainingToLocalStorage(training);
  };

  return (
    <Switch
      fallback={
        <>
          <DismissibleModal
            open={openConfirmationModal()}
            setOpen={setOpenConfirmationModal}
            message={t('confirm.training.delete.message')}
            confirmLabel={t('confirm.delete.training')}
            cancelLabel={t('no.cancel')}
            onConfirm={() => {
              clearTrainingFromLocalStorage();
              navigate('/', { replace: true });
            }}
          />
          <TrainingPreview
            training={training}
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
                  onCreateSet(cloneSet(training.sets[setIdx]));
                  saveTrainingToLocalStorage(training);
                },
              },
              {
                text: t('move.up'),
                icon: 'fa-arrow-up',
                onClick: (setIdx) => {
                  moveSetUpInNewTraining(setIdx, setTraining);
                  saveTrainingToLocalStorage(training);
                },
                disabledFunc: (setIdx) => setIdx === 0,
              },
              {
                text: t('move.down'),
                icon: 'fa-arrow-down',
                onClick: (setIdx) => {
                  moveSetDownInNewTraining(
                    setIdx,
                    training.sets.length,
                    setTraining
                  );
                  saveTrainingToLocalStorage(training);
                },
                disabledFunc: (setIdx) => setIdx === training.sets.length - 1,
              },
              {
                text: t('delete'),
                icon: 'fa-trash',
                onClick: (setIdx) => {
                  deleteSetInNewTraining(setIdx, setTraining);
                  saveTrainingToLocalStorage(training);
                },
              },
            ]}
            rightHeaderComponent={() => (
              <button
                class="rounded-lg p-1 text-right"
                onClick={() => setOpenConfirmationModal(true)}
              >
                <i class="fa-trash fa-solid fa-xl cursor-pointer text-red-500"></i>
              </button>
            )}
          />

          <TrainingEditButtonGroup
            training={training}
            backLabel={t('back')}
            onBack={() => navigate('/', { replace: true })}
            onAddSet={() => setShowCreateSet(true)}
            confirmLabel={t('finish')}
            onConfirm={() => setShowTrainingSession(true)}
          />
        </>
      }
    >
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
          onSubmit={onTrainingSessionSubmit}
          onBack={() => setShowTrainingSession(false)}
        />
      </Match>
    </Switch>
  );
};

export default NewTrainingPage;
