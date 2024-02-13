import { useTransContext } from '@mbarzda/solid-i18next'
import {
  Component,
  createSignal,
  Match,
  onCleanup,
  onMount,
  Show,
  Switch,
} from 'solid-js'
import { createStore } from 'solid-js/store'
import ConfirmationModal from '../components/ConfirmationModal'
import { clearOnBackOverride, setOnBackOverride } from '../components/Header'
import {
  NewTraining,
  NewTrainingSet,
  Training,
  TrainingSet,
} from 'swimlogs-api'
import { cloneSet } from '../lib/clone'
import {
  clearTrainingFromLocalStorage,
  saveTrainingToLocalStorage,
} from '../state/local-storage'
import EditSetPage from './EditSetPage'
import EditTrainingSessionPage from './EditTrainingSessionPage'
import TrainingPreviewPage from './TrainingPreviewPage'

interface EditTrainingPageProps {
  training: ReturnType<typeof createStore<NewTraining | Training>>
  onSubmit: (training: NewTraining | Training) => void
  onDelete?: () => void
  onBack: () => void
  saveToLocalStorage: boolean
}

const EditTrainingPage: Component<EditTrainingPageProps> = (props) => {
  const [t] = useTransContext()
  const [showCreateSet, setShowCreateSet] = createSignal(false)
  const [editedSetIdx, setEditedSetIdx] = createSignal(-1)
  const [showTrainingSession, setShowTrainingSession] = createSignal(false)
  const [training, setTraining] = props.training

  onMount(() => setOnBackOverride(onBack))
  onCleanup(() => clearOnBackOverride())

  const onBack = () => {
    if (showCreateSet()) {
      setShowCreateSet(false)
    } else if (editedSetIdx() !== -1) {
      setEditedSetIdx(-1)
    } else if (showTrainingSession()) {
      setShowTrainingSession(false)
    } else {
      props.onBack()
    }
  }

  const onDeleteTraining = () => {
    if (props.saveToLocalStorage) {
      clearTrainingFromLocalStorage()
    }
    props.onDelete?.()
  }

  const onTrainingSessionSubmit = (trainingSession: {
    start: Date
    durationMin: number
  }) => {
    setShowTrainingSession(false)
    setTraining('start', trainingSession.start)
    setTraining('durationMin', trainingSession.durationMin)
    props.onSubmit(training)
    if (props.saveToLocalStorage) {
      clearTrainingFromLocalStorage()
    }
  }

  const recalculateTotalDistance = () => {
    let newTotalDistance = 0
    training.sets.forEach(
      (s) => (newTotalDistance += s.repeat * (s.distanceMeters ?? 0))
    )
    setTraining('totalDistance', newTotalDistance)
  }

  const onCreateSet = (set: NewTrainingSet) => {
    setShowCreateSet(false)
    set.setOrder = training.sets.length
    setTraining('sets', [...training.sets, set])
    recalculateTotalDistance()
    if (props.saveToLocalStorage) {
      saveTrainingToLocalStorage(training)
    }
  }

  const onEditSet = (set: NewTrainingSet | TrainingSet) => {
    setTraining('sets', (sets) => {
      const tmp = sets[editedSetIdx()]
      set.setOrder = tmp.setOrder
      if ('id' in tmp) {
        // @ts-ignore
        set['id'] = tmp.id
      }
      sets[editedSetIdx()] = set
      return sets
    })
    setEditedSetIdx(-1)
    recalculateTotalDistance()
    if (props.saveToLocalStorage) {
      saveTrainingToLocalStorage(training)
    }
  }

  const onMoveUpSet = (setIdx: number) => {
    if (setIdx === 0) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx - 1]
      sets[setIdx - 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
    if (props.saveToLocalStorage) {
      saveTrainingToLocalStorage(training)
    }
  }

  const onMoveDownSet = (setIdx: number) => {
    if (setIdx === training.sets.length - 1) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx + 1]
      sets[setIdx + 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
    if (props.saveToLocalStorage) {
      saveTrainingToLocalStorage(training)
    }
  }

  const onDeleteSet = (setIdx: number) => {
    setTraining(
      'totalDistance',
      training.totalDistance - training.sets[setIdx].totalDistance
    )
    setTraining('sets', (sets) =>
      sets.filter((_, i) => i !== setIdx).map((s, i) => ({ ...s, setOrder: i }))
    )
    if (props.saveToLocalStorage) {
      saveTrainingToLocalStorage(training)
    }
  }

  return (
    <div class="pb-4">
      <Switch>
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
            onBack={onBack}
          />
        </Match>
        <Match when={!showCreateSet()}>
          <TrainingPreviewPage
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
                  onCreateSet(cloneSet(training.sets[setIdx]))
                  if (props.saveToLocalStorage) {
                    saveTrainingToLocalStorage(training)
                  }
                },
              },
              {
                text: t('move.up'),
                icon: 'fa-arrow-up',
                onClick: onMoveUpSet,
                disabledFunc: (setIdx) => setIdx === 0,
              },
              {
                text: t('move.down'),
                icon: 'fa-arrow-down',
                onClick: onMoveDownSet,
                disabledFunc: (setIdx) => setIdx === training.sets.length - 1,
              },
              { text: t('delete'), icon: 'fa-trash', onClick: onDeleteSet },
            ]}
            rightHeaderComponent={() => (
              <div class="text-right">
                <ConfirmationModal
                  icon="fa-trash"
                  iconColor="text-red-500"
                  message={t('confirm.training.delete.message')}
                  confirmLabel={t('confirm.delete.training')}
                  cancelLabel={t('no.cancel')}
                  onConfirm={onDeleteTraining}
                  onCancel={() => {}}
                />
              </div>
            )}
          />

          <div class="py-2 text-center">
            <button
              class="h-12 w-12 rounded-full bg-sky-500"
              onClick={() => setShowCreateSet(true)}
            >
              <i class="fa-solid fa-plus fa-2xl text-white"></i>
            </button>
          </div>

          <Show when={training.sets.length > 0}>
            <div class="flex items-center justify-between p-4 md:justify-around">
              <button
                class="w-24 rounded-lg bg-red-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-red-300"
                onClick={onBack}
              >
                {t('back')}
              </button>

              <button
                class="w-24 rounded-lg bg-green-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-green-300"
                onClick={() => setShowTrainingSession(true)}
              >
                {t('finish')}
              </button>
            </div>
          </Show>
        </Match>
      </Switch>
    </div>
  )
}

export default EditTrainingPage
