import { Trans, useTransContext } from '@mbarzda/solid-i18next';
import { useNavigate } from '@solidjs/router';
import { createStore, SetStoreFunction, unwrap } from 'solid-js/store';
import { showToast } from '../App';
import { ToastMode } from '../components/DismissibleToast';
import { NewTraining, NewTrainingSet, StartTypeEnum } from 'swimlogs-api';
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
  Show,
  Switch,
} from 'solid-js';
import {
  deleteSetInNewTraining,
  moveSetDownInNewTraining,
  moveSetUpInNewTraining,
  recalculateTotalDistance,
} from '../lib/training';
import TrainingPreview from '../components/TrainingPreview';
import { cloneSet } from '../lib/clone';
import DismissibleModal from '../components/DismissibleModal';
import { useOnBackcontext } from './Routing';
import { clearHeaderButton, setHeaderButton } from '../components/Header';
import SetEditForm from '../components/SetEditForm';
import { SmallIntMax } from '../lib/consts';
import { Callout, CalloutTitle } from '../components/ui/callout';
import SessionEditForm from '../components/SessionEditForm';

type PreviewScreen = { screen: 'preview' };
type CreateSetScreen = { screen: 'create-set' };
type EditSetScreen = {
  screen: 'edit-set';
  setIdx: number;
  setStore: [get: NewTrainingSet, set: SetStoreFunction<NewTrainingSet>];
};
type SessionScreen = { screen: 'session' };

type PageOnScreen =
  | PreviewScreen
  | CreateSetScreen
  | EditSetScreen
  | SessionScreen;

const NewTrainingPage: Component = () => {
  const [t] = useTransContext();
  const navigate = useNavigate();
  const [onBack, setOnBack] = useOnBackcontext();

  const [openConfirmationModal, setOpenConfirmationModal] = createSignal(false);
  const [training, setTraining] = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: new Date(new Date().setHours(new Date().getHours(), 0, 0, 0)),
      durationMin: 60,
      totalDistance: 0,
      sets: [],
    }
  );

  const [screen, setScreen] = createStore<PageOnScreen>({ screen: 'preview' });
  const [newSet, setNewSet] = createStore<NewTrainingSet>({
    setOrder: training.sets.length,
    repeat: 1,
    distanceMeters: 100,
    startType: StartTypeEnum.None,
    totalDistance: 100,
  });

  const isSetValid = (set: NewTrainingSet) => {
    const isStartValid =
      set.startType === StartTypeEnum.None ||
      (set.startSeconds && set.startSeconds !== 0);
    const isDistanceValid =
      set.distanceMeters >= 25 && set.distanceMeters <= SmallIntMax;

    return isDistanceValid && isStartValid;
  };

  function onBackOverride() {
    if (screen.screen === 'create-set') {
      setScreen({ screen: 'preview' });
    } else if (screen.screen === 'edit-set') {
      setScreen({ screen: 'preview' });
    } else if (screen.screen === 'session') {
      setScreen({ screen: 'preview' });
    } else {
      setOnBack();
      onBack();
    }
  }

  const headerButton = () => {
    return (
      <i
        classList={{
          'fa-arrow-right fa-2xl': screen.screen === 'preview',
          'fa-check fa-2xl ': screen.screen !== 'preview',
          'text-white/50 pointer-events-none':
            (screen.screen === 'create-set' && !isSetValid(newSet)) ||
            (screen.screen === 'edit-set' && !isSetValid(screen.setStore[0])),
        }}
        class="text-right fa-solid cursor-pointer text-white"
        onClick={() => {
          switch (screen.screen) {
            case 'session':
              // TODO set the session and send it
              // onTrainingSessionSubmit
              break;
            case 'create-set':
              submitNewSet();
              setScreen({ screen: 'preview' });
              break;
            case 'edit-set':
              editSet();
              setScreen({ screen: 'preview' });
              break;
            case 'preview':
              setScreen({ screen: 'session' });
              break;
          }
        }}
      />
    );
  };

  onMount(() => {
    setOnBack(onBackOverride);
    setHeaderButton(headerButton());
  });

  onCleanup(() => {
    setOnBack();
    clearHeaderButton();
  });

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
    setScreen({ screen: 'preview' });
    setTraining('start', trainingSession.start);
    setTraining('durationMin', trainingSession.durationMin);
    onSubmit(training);
    clearTrainingFromLocalStorage();
  };

  const submitNewSet = () => {
    const set: NewTrainingSet = {
      setOrder: newSet.setOrder,
      repeat: newSet.repeat,
      distanceMeters: newSet.distanceMeters,
      description: newSet.description,
      startType: newSet.startType,
      startSeconds: newSet.startSeconds,
      totalDistance: newSet.repeat * newSet.distanceMeters,
      equipment:
        (newSet.equipment ?? []).length > 0 ? newSet.equipment : undefined,
      group: newSet.group,
    };

    onCreateSet(set);
  };

  const onCreateSet = (set: NewTrainingSet) => {
    setScreen({ screen: 'preview' });
    set.setOrder = training.sets.length;
    setTraining('sets', [...training.sets, set]);
    setTraining('totalDistance', recalculateTotalDistance(training));
    saveTrainingToLocalStorage(training);
  };

  const editSet = () => {
    setTraining('sets', (sets) => {
      if (screen.screen !== 'edit-set') {
        console.error(
          `screen is not in edit-set state, but is in ${screen.screen}`
        );
        return sets;
      }

      const [set, setSet] = screen.setStore;
      const idx = screen.setIdx;
      const tmp = sets[idx];
      setSet('setOrder', tmp.setOrder);
      sets[idx] = set;
      return sets;
    });

    setScreen({ screen: 'preview' });
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
                onClick: (setIdx) =>
                  setScreen({
                    screen: 'edit-set',
                    setIdx,
                    setStore: createStore(unwrap(training.sets[setIdx])),
                  }),
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
          />

          <div class="p-4">
            <Show
              when={training.sets.length !== 0}
              fallback={
                <Callout>
                  <CalloutTitle class="whitespace-pre-line">
                    <Trans key="no.sets.in.training" />
                  </CalloutTitle>
                </Callout>
              }
            >
              <button onClick={() => setOpenConfirmationModal(true)}>
                <i class="fa-solid fa-trash-can text-red-500 fa-2xl"></i>
              </button>
            </Show>
          </div>

          <button
            class="fixed bottom-4 right-4 h-12 w-12 rounded-full bg-sky-500"
            onClick={() => setScreen({ screen: 'create-set' })}
          >
            <i class="fa-solid fa-plus fa-2xl text-white"></i>
          </button>
        </>
      }
    >
      <Match when={screen.screen === 'create-set'}>
        <SetEditForm set={newSet} updateSet={setNewSet} />
      </Match>
      <Match when={screen.screen === 'edit-set'}>
        <SetEditForm
          set={(screen as EditSetScreen).setStore[0]}
          updateSet={(screen as EditSetScreen).setStore[1]}
        />
      </Match>
      <Match when={screen.screen === 'session'}>
        <SessionEditForm training={training} updateTraining={setTraining} />
      </Match>
    </Switch>
  );
};

export default NewTrainingPage;
