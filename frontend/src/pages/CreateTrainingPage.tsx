import { useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, Match, Show, Switch } from 'solid-js'
import { createStore } from 'solid-js/store'
import { setOnBackOverrideOnce } from '../components/Header'
import { NewTraining, NewTrainingSet } from '../generated'
import { cloneSet } from '../lib/clone'
import {
  clearTrainingFromLocalStorage,
  loadTrainingFromLocalStorage,
  saveTrainingToLocalStorage,
} from '../state/local-storage'
import EditSetPage from './EditSetPage'
import EditTrainingSessionPage from './EditTrainingSessionPage'
import TrainingPreviewPage from './TrainingPreviewPage'

const CreateTrainingPage: Component = () => {
  const [t] = useTransContext()
  const [showCreateSet, setShowCreateSet] = createSignal(false)
  const [editedSetIdx, setEditedSetIdx] = createSignal(-1)
  const [showTrainingSession, setShowTrainingSession] = createSignal(true)

  const [training, setTraining] = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: new Date(),
      durationMin: 0,
      totalDistance: 0,
      sets: [],
    }
  )

  const onBack = () => {
    if (showCreateSet()) {
      setShowCreateSet(false)
    } else if (editedSetIdx() !== -1) {
      setEditedSetIdx(-1)
    } else if (showTrainingSession()) {
      setShowTrainingSession(false)
    } else {
      history.back()
    }
  }

  const onDeleteTraining = () => {
    clearTrainingFromLocalStorage()
    history.back()
  }

  const onTrainingSessionSubmit = (trainingSession: NewTraining) => {
	setShowTrainingSession(false)
    console.log('onTrainingSessionSubmit', trainingSession)
  }

  const onCreateSet = (set: NewTrainingSet) => {
    setShowCreateSet(false)
    set.setOrder = training.sets.length
    setTraining('sets', [...training.sets, set])
    setTraining('totalDistance', training.totalDistance + set.totalDistance)
    saveTrainingToLocalStorage(training)
  }

  const onEditSet = (set: NewTrainingSet) => {
    setTraining('sets', (sets) => {
      const tmp = sets[editedSetIdx()]
      set.setOrder = tmp.setOrder
      sets[editedSetIdx()] = set
      return sets
    })
    setEditedSetIdx(-1)
    saveTrainingToLocalStorage(training)
  }

  const onMoveUpSet = (setIdx: number) => {
    if (setIdx === 0) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx - 1]
      sets[setIdx - 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
    saveTrainingToLocalStorage(training)
  }

  const onMoveDownSet = (setIdx: number) => {
    if (setIdx === training.sets.length - 1) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx + 1]
      sets[setIdx + 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
    saveTrainingToLocalStorage(training)
  }

  const onDeleteSet = (setIdx: number) => {
    setTraining(
      'totalDistance',
      training.totalDistance - training.sets[setIdx].totalDistance
    )
    setTraining('sets', (sets) =>
      sets.filter((_, i) => i !== setIdx).map((s, i) => ({ ...s, setOrder: i }))
    )
    saveTrainingToLocalStorage(training)
  }

  return (
    <div class="pb-24 md:pb-4">
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
          <EditTrainingSessionPage onSubmit={onTrainingSessionSubmit} />
        </Match>
        <Match when={!showCreateSet()}>
          <TrainingPreviewPage
            training={training}
            showOptions={true}
            options={{
              onEdit: (setIdx) => {
                setOnBackOverrideOnce(onBack)
                setEditedSetIdx(setIdx)
              },
              onDuplicate: (setIdx) => {
                onCreateSet(cloneSet(training.sets[setIdx]))
                saveTrainingToLocalStorage(training)
              },
              onMoveUp: onMoveUpSet,
              onMoveDown: onMoveDownSet,
              onDelete: onDeleteSet,
            }}
            showDeleteTraining={true}
            onDeleteTraining={onDeleteTraining}
          />

          <Show when={training.sets.length > 0}>
            <div class="flex items-center justify-between p-4 md:justify-around">
              <button
                class="w-24 rounded-lg bg-red-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-red-300"
                onClick={() => onBack()}
              >
                {t('back')}
              </button>

              <button
                class="w-24 rounded-lg bg-green-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-green-300"
                onClick={() => {
                  setOnBackOverrideOnce(onBack)
                  setShowTrainingSession(true)
                }}
              >
                {t('create')}
              </button>
            </div>
          </Show>

          <button
            class="fixed bottom-2 right-2 h-16 w-16 rounded-lg bg-sky-500"
            onClick={() => {
              setOnBackOverrideOnce(onBack)
              setShowCreateSet(true)
            }}
          >
            <i class="fa-solid fa-plus fa-2xl text-white"></i>
          </button>
        </Match>
      </Switch>
    </div>
  )
}

export default CreateTrainingPage
