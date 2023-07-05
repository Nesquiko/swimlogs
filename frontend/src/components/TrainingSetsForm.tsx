import { Trans } from '@mbarzda/solid-i18next'
import { batch, Component, For, Show } from 'solid-js'
import { produce } from 'solid-js/store'
import { NewTrainingSet, StartType } from '../generated'
import { cloneSet } from '../lib/clone'
import { isInvalidTrainingEmpty, validateTraining } from '../lib/validation'
import { useCreateTraining } from './context/CreateTrainingContextProvider'
import SetForm from './SetForm'

const TrainingSetsForm: Component = () => {
  const [
    { training, setTraining },
    { invalidTraining, setInvalidTraining },
    ,
    ,
    [, setCurrentComponent]
  ] = useCreateTraining()

  const sumbit = () => {
    setInvalidTraining(validateTraining(training))

    if (!isInvalidTrainingEmpty(invalidTraining)) {
      return
    }

    setCurrentComponent((c) => c + 1)
  }

  const addNewSet = () => {
    const newSet = {
      repeat: 1,
      setOrder: training.sets.length,
      distanceMeters: 100,
      startType: StartType.None,
      totalDistance: 100
    } as NewTrainingSet
    addSet(newSet)
  }

  const duplicateSet = (idx: number) => {
    const newSet = cloneSet(training.sets[idx])
    addSet(newSet)
  }

  const addSet = (s: NewTrainingSet) => {
    batch(() => {
      setTraining(
        'sets',
        produce((sets) => sets.push(s))
      )
      setInvalidTraining(
        'invalidSets',
        produce((sets) => sets?.push({}))
      )
    })
  }

  const deleteSet = (idx: number) => {
    batch(() => {
      setTraining('sets', (sets) => sets.filter((_, i) => i !== idx))
      setInvalidTraining('invalidSets', (sets) =>
        sets?.filter((_, i) => i !== idx)
      )
    })
  }

  return (
    <div class="m-4">
      <p class="my-4 text-xl">
        <Trans key="total.distance.training" />{' '}
        {training.totalDistance.toLocaleString()}m
      </p>
      <Show
        when={training.sets.length !== 0}
        fallback={
          <div class="m-4 flex items-center justify-start rounded bg-blue-200 p-4 text-xl font-bold">
            <Trans key="no.sets.in.training" />
          </div>
        }
      >
        <For each={training.sets}>
          {(set, setIdx) => {
            return (
              <SetForm
                set={set}
                setIdx={setIdx()}
                invalidSet={invalidTraining.invalidSets![setIdx()]}
                onDelete={() => deleteSet(setIdx())}
                onDuplicate={() => duplicateSet(setIdx())}
              />
            )
          }}
        </For>
      </Show>
      <button
        class="float-right my-4 rounded bg-sky-500 p-2 font-bold text-white"
        onClick={() => addNewSet()}
      >
        <Trans key="add.set" />
      </button>

      <button
        class="fixed bottom-0 right-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
        onClick={() => sumbit()}
      >
        <Trans key="next" />
      </button>
      <button
        class="fixed bottom-0 left-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
        onClick={() => setCurrentComponent((c) => c - 1)}
      >
        <Trans key="previous" />
      </button>
    </div>
  )
}

export default TrainingSetsForm
