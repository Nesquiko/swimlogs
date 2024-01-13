import { Component, createSignal, Match, Switch } from 'solid-js'
import { createStore } from 'solid-js/store'
import { Equipment, NewTraining, NewTrainingSet } from '../generated'
import { cloneSet } from '../lib/clone'
import CreateSetPage from './CreateSetPage'
import TrainingPreviewPage from './TrainingPreviewPage'

const CreateTrainingPage: Component = () => {
  const [showCreateSet, setShowCreateSet] = createSignal(false)

  const [training, setTraining] = createStore<NewTraining>({
    start: new Date(),
    durationMin: 0,
    totalDistance: 1200,
    sets: [
      {
        setOrder: 0,
        repeat: 8,
        distanceMeters: 50,
        startType: 'Pause',
        startSeconds: 45,
        totalDistance: 400,
        equipment: [Equipment.Fins, Equipment.Paddles, Equipment.Monofin],
      },
      {
        setOrder: 1,
        repeat: 1,
        distanceMeters: 400,
        startType: 'None',
        startSeconds: 0,
        totalDistance: 400,
        description:
          'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation',
        equipment: [Equipment.Paddles],
      },
      {
        setOrder: 2,
        repeat: 4,
        distanceMeters: 100,
        startType: 'Interval',
        startSeconds: 90,
        totalDistance: 400,
        description: 'Kraul',
        equipment: [
          Equipment.Fins,
          Equipment.Board,
          Equipment.Paddles,
          Equipment.Snorkel,
          Equipment.Monofin,
        ],
      },
    ],
  })

  const onCreateSet = (set: NewTrainingSet) => {
    setShowCreateSet(false)
    set.setOrder = training.sets.length
    setTraining('sets', [...training.sets, set])
    setTraining('totalDistance', training.totalDistance + set.totalDistance)
  }

  const onEditSet = (setIdx: number) => {}

  const onDuplicateSet = (setIdx: number) => {
    onCreateSet(cloneSet(training.sets[setIdx]))
  }

  const onMoveUpSet = (setIdx: number) => {
    if (setIdx === 0) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx - 1]
      sets[setIdx - 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
  }

  const onMoveDownSet = (setIdx: number) => {
    if (setIdx === training.sets.length - 1) return

    setTraining('sets', (sets) => {
      const tmp = sets[setIdx]
      sets[setIdx] = sets[setIdx + 1]
      sets[setIdx + 1] = tmp
      return sets.map((s, i) => ({ ...s, setOrder: i }))
    })
  }

  const onDeleteSet = (setIdx: number) => {
    setTraining(
      'totalDistance',
      training.totalDistance - training.sets[setIdx].totalDistance
    )
    setTraining('sets', (sets) =>
      sets.filter((_, i) => i !== setIdx).map((s, i) => ({ ...s, setOrder: i }))
    )
  }

  return (
    <div>
      <Switch>
        <Match when={showCreateSet()}>
          <CreateSetPage onCreateSet={onCreateSet} />
        </Match>
        <Match when={!showCreateSet()}>
          <TrainingPreviewPage
            training={training}
            showOptions={true}
            options={{
              onEdit: onEditSet,
              onDuplicate: onDuplicateSet,
              onMoveUp: onMoveUpSet,
              onMoveDown: onMoveDownSet,
              onDelete: onDeleteSet,
            }}
          />

          <button
            class="fixed bottom-2 right-2 h-16 w-16 rounded-lg bg-sky-500"
            onClick={() => setShowCreateSet(true)}
          >
            <i class="fa-solid fa-plus fa-2xl text-white"></i>
          </button>
        </Match>
      </Switch>
    </div>
  )
}

export default CreateTrainingPage
