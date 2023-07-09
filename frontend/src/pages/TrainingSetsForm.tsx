import { Trans } from '@mbarzda/solid-i18next'
import { t } from 'i18next'
import { Component, createSignal, For, Match, Show, Switch } from 'solid-js'
import { produce } from 'solid-js/store'
import { cloneSet } from '../lib/clone'
import plusSvg from '../assets/plus.svg'
import { useCreateTraining } from '../components/CreateTrainingContextProvider'
import { NewTrainingSet, StartType } from '../generated'
import SetModal from '../components/SetModal'
import MenuModal from '../components/MenuModal'
import SuperSetEditPage from './SuperSetEditPage'
import SetCard from '../components/SetCard'

const TrainingSetsForm: Component = () => {
  const [{ training, setTraining }, {}, , , [, setCurrentComponent]] =
    useCreateTraining()

  const newSet = (): NewTrainingSet => {
    return {
      repeat: 1,
      distanceMeters: 100,
      startType: StartType.None,
      totalDistance: 100
    }
  }

  const newSuperSet = (): NewTrainingSet => {
    return {
      repeat: 1,
      distanceMeters: 100,
      startType: StartType.None,
      totalDistance: 100,
      subSets: []
    }
  }

  const [addMenuModalOpener, setAddMenuModalOpener] = createSignal({})
  const [setSettingsOpener, setSetSettingsOpener] = createSignal({ idx: -1 })

  const [addModalOpener, setAddModalOpener] = createSignal({
    set: newSet()
  })
  const [editModalOpener, setEditModalOpener] = createSignal({
    set: newSet(),
    idx: -1
  })
  const [superSetFormOpen, setSuperSetFormOpen] = createSignal(false)

  const addNewSet = (set: NewTrainingSet) => {
    set.setOrder = training.sets.length
    set.subSets?.forEach((s) => {
      s.setOrder = set.setOrder
    })
    addSet(set)
  }

  const duplicateSet = (idx: number) => {
    const newSet = cloneSet(training.sets[idx])
    addSet(newSet)
  }

  const addSet = (s: NewTrainingSet) => {
    s.setOrder = training.sets.length
    setTraining(
      'sets',
      produce((sets) => sets.push(s))
    )
  }

  const editSet = (set: NewTrainingSet, idx: number) => {
    setTraining(
      'sets',
      produce((sets) => (sets[idx] = set))
    )
  }

  const deleteSet = (idx: number) => {
    setTraining('sets', (sets) =>
      sets.filter((_, i) => i !== idx).map((s, i) => ({ ...s, setOrder: i }))
    )
  }

  return (
    <div class="m-4">
      <Switch>
        <Match when={superSetFormOpen()}>
          <SuperSetEditPage
            superSet={newSuperSet()}
            onSubmitSuperSet={(set) => {
              addNewSet(set)
              setSuperSetFormOpen(false)
            }}
            onClose={() => setSuperSetFormOpen(false)}
          />
        </Match>
        <Match when={!superSetFormOpen()}>
          <SetModal opener={addModalOpener()} onSubmitSet={addNewSet} />
          <SetModal
            opener={editModalOpener()}
            onSubmitSet={(set, idx) => editSet(set, idx!)}
            submitBtnLabelKey="edit"
          />
          <MenuModal
            widthRem="15"
            opener={addMenuModalOpener()}
            items={[
              {
                label: t('add.new.set', 'Add new set'),
                action: () => setAddModalOpener({ set: newSet() })
              },
              {
                label: t('add.new.super.set', 'Add new superset'),
                action: () => setSuperSetFormOpen(true)
              }
            ]}
          />
          <MenuModal<{ idx: number }>
            opener={setSettingsOpener()}
            items={[
              {
                label: t('edit', 'Edit'),
                action: (o) =>
                  setEditModalOpener({
                    set: training.sets[o.idx],
                    idx: o.idx
                  })
              },
              {
                label: t('duplicate', 'Duplicate'),
                action: (o) => duplicateSet(o.idx)
              },
              {
                label: t('delete', 'Delete'),
                action: (o) => deleteSet(o.idx)
              }
            ]}
            header={(o) => t('set', 'Set') + ' ' + (o.idx + 1)}
          />

          <Show
            when={training.sets.length !== 0}
            fallback={
              <div class="m-4 rounded-lg bg-sky-200 p-4 text-xl font-semibold">
                <Trans key="no.sets.in.training" />
              </div>
            }
          >
            <For each={training.sets}>
              {(set, setIdx) => {
                return (
                  <SetCard
                    set={set}
                    setNumber={setIdx() + 1}
                    onSettingsClick={() =>
                      setSetSettingsOpener({ idx: setIdx() })
                    }
                  />
                )
              }}
            </For>
          </Show>
          <p class="my-4 text-xl">
            <Trans key="total.distance.training" />{' '}
            {training.totalDistance.toLocaleString()}m
          </p>
          <button
            class="float-right h-10 w-10 rounded-full bg-green-500 text-2xl text-white shadow"
            onClick={() => setAddMenuModalOpener({})}
          >
            <img src={plusSvg} />
          </button>
          <div class="h-32 w-full"></div>

          <button
            class="fixed bottom-0 right-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
            onClick={() => setCurrentComponent((c) => c + 1)}
          >
            <Trans key="next" />
          </button>
          <button
            class="fixed bottom-0 left-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
            onClick={() => setCurrentComponent((c) => c - 1)}
          >
            <Trans key="previous" />
          </button>
        </Match>
      </Switch>
    </div>
  )
}

export default TrainingSetsForm
