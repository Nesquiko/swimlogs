import { Trans } from '@mbarzda/solid-i18next'
import { t } from 'i18next'
import { Component, createSignal, For, Match, Show, Switch } from 'solid-js'
import { produce } from 'solid-js/store'
import { cloneSet } from '../lib/clone'
import plusSvg from '../assets/plus.svg'
import { useCreateTraining } from '../components/CreateTrainingContextProvider'
import { NewTrainingSet } from '../generated'
import SetModal from '../components/SetModal'
import MenuModal from '../components/MenuModal'
import SuperSetEditPage from './SuperSetEditPage'
import SetCard from '../components/SetCard'

const TrainingSetsForm: Component = () => {
  const [{ training, setTraining }, {}, , , [, setCurrentComponent]] =
    useCreateTraining()

  const [addMenuModalOpener, setAddMenuModalOpener] = createSignal({})
  const [setSettingsOpener, setSetSettingOpener] = createSignal({ idx: -1 })

  const [setModalOpen, setSetModalOpen] = createSignal({})
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
    setTraining(
      'sets',
      produce((sets) => sets.push(s))
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
            onAddSet={(set) => {
              addNewSet(set)
              setSuperSetFormOpen(false)
            }}
            onClose={() => setSuperSetFormOpen(false)}
          />
        </Match>
        <Match when={!superSetFormOpen()}>
          <SetModal open={setModalOpen()} onAddSet={addNewSet} />
          <MenuModal
            widthRem="15"
            opener={addMenuModalOpener()}
            items={[
              {
                label: t('add.new.set', 'Add new set'),
                action: () => setSetModalOpen({})
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
                    onSettingsClick={() =>
                      setSetSettingOpener({ idx: setIdx() })
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
