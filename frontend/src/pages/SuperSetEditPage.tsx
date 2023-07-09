import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, For, Show } from 'solid-js'
import { createStore, produce, unwrap } from 'solid-js/store'
import MenuModal from '../components/MenuModal'
import SetCard from '../components/SetCard'
import SetModal from '../components/SetModal'
import { openToast, ToastType } from '../components/Toast'
import { NewTrainingSet, StartType } from '../generated'
import { cloneSet } from '../lib/clone'
import { SmallIntMax } from '../lib/consts'

type SuperSetEditPageProps = {
  superSet: NewTrainingSet
  onSubmitSuperSet: (set: NewTrainingSet) => void
  onClose: () => void
}

const SuperSetEditPage: Component<SuperSetEditPageProps> = (props) => {
  const [superSet, setSuperSet] = createStore<NewTrainingSet>(props.superSet)

  const newSubSet = (): NewTrainingSet => {
    return {
      repeat: 1,
      distanceMeters: 100,
      startType: StartType.None,
      totalDistance: 100
    }
  }

  const [setModalOpener, setSetModalOpener] = createSignal({
    set: newSubSet()
  })
  const [setSettingsOpener, setSetSettingsOpener] = createSignal({ idx: -1 })
  const [editModalOpener, setEditModalOpener] = createSignal(
    {
      set: newSubSet(),
      idx: -1
    },
    { equals: false }
  )

  const [t] = useTransContext()

  const addNewSubSet = (set: NewTrainingSet) => {
    set.setOrder = superSet.setOrder
    set.subSetOrder = superSet.subSets?.length

    setSuperSet(
      'subSets',
      produce((subSets) => subSets?.push(set))
    )
  }

  const editSubSet = (subSet: NewTrainingSet, idx: number) => {
    setSuperSet(
      'subSets',
      produce((subSets) => (subSets![idx] = subSet))
    )
  }

  const duplicateSubSet = (idx: number) => {
    const newSet = cloneSet(superSet.subSets![idx])
    newSet.subSetOrder = superSet.subSets?.length
    addNewSubSet(newSet)
  }

  const deleteSet = (idx: number) => {
    setSuperSet('subSets', (subsets) =>
      subsets
        ?.filter((_, i) => i !== idx)
        .map((s, i) => ({ ...s, subSetOrder: i }))
    )
  }

  const isSetValid = (): boolean => {
    if (superSet.repeat < 1) {
      return false
    }

    if (
      superSet.startType !== StartType.None &&
      superSet.startSeconds !== undefined &&
      superSet.startSeconds < 1
    ) {
      return false
    }

    if (superSet.subSets?.length === 0) {
      openToast(
        t('add.at.least.one.set', 'Add at least one set'),
        ToastType.ERROR
      )
      return false
    }

    return true
  }

  return (
    <div class="h-full w-full rounded-lg">
      <SetModal opener={setModalOpener()} onSubmitSet={addNewSubSet} />
      <SetModal
        opener={editModalOpener()}
        onSubmitSet={(set, idx) => editSubSet(set, idx!)}
        submitBtnLabelKey="edit"
      />
      <MenuModal<{ idx: number }>
        opener={setSettingsOpener()}
        items={[
          {
            label: t('edit', 'Edit'),
            action: (o) => {
              setEditModalOpener({
                set: superSet.subSets![o.idx],
                idx: o.idx
              })
            }
          },
          {
            label: t('duplicate', 'Duplicate'),
            action: (o) => duplicateSubSet(o.idx)
          },
          {
            label: t('delete', 'Delete'),
            action: (o) => deleteSet(o.idx)
          }
        ]}
        header={(o) => t('set', 'Set') + ' ' + (o.idx + 1)}
      />

      <p class="text-center text-2xl">
        <Trans key="add.new.super.set" />
      </p>
      <hr class="my-2 rounded-lg border-2 border-slate-500" />
      <div class="my-2 flex items-center justify-between">
        <label class="text-xl" for="repeat">
          <Trans key="repeat" />
        </label>
        <input
          id="repeat"
          type="number"
          placeholder="1"
          classList={{
            'border-red-500 text-red-500': superSet.repeat < 1,
            'border-slate-300': superSet.repeat >= 1
          }}
          class="w-24 rounded-md border p-2 text-center text-lg focus:border-blue-500 focus:outline-none focus:ring"
          value={superSet.repeat}
          onChange={(e) => {
            let repeat = parseInt(e.target.value)
            if (Number.isNaN(repeat) || repeat < 1 || repeat > SmallIntMax) {
              repeat = 0
            }
            setSuperSet('repeat', repeat)

            const subSetsDist =
              superSet.subSets?.reduce((acc, subSet) => {
                return acc + subSet.totalDistance
              }, 0) || 0

            setSuperSet('totalDistance', repeat * subSetsDist)
          }}
        />
      </div>
      <div class="my-2 flex items-center justify-between">
        <label class="text-xl" for="start">
          <Trans key="start" />
        </label>
        <select
          id="start"
          class="w-32 rounded-md border border-solid border-slate-300 bg-white p-2 text-center text-lg focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            setSuperSet('startType', e.target.value as StartType)
          }}
        >
          <For each={Object.keys(StartType)}>
            {(typ) => (
              <option selected={typ === superSet.startType} value={typ}>
                <Trans key={typ.toLowerCase()} />
              </option>
            )}
          </For>
        </select>
      </div>
      <div
        classList={{
          visible: superSet.startType !== StartType.None,
          invisible: superSet.startType === StartType.None
        }}
        class="my-2 flex items-center justify-between"
      >
        <label class="text-xl" for="seconds">
          <Trans key="seconds" />
        </label>
        <input
          id="seconds"
          type="number"
          placeholder="20"
          classList={{
            'border-red-500 text-red-500':
              superSet.startSeconds !== undefined && superSet.startSeconds < 1,
            'border-slate-300':
              superSet.startSeconds === undefined || superSet.startSeconds >= 1
          }}
          class="w-24 rounded-md border border-solid border-slate-300 bg-white p-2 text-center text-lg focus:border-sky-500 focus:outline-none focus:ring"
          value={superSet.startSeconds}
          onChange={(e) => {
            const val = e.target.value
            let seconds = parseInt(val)
            if (Number.isNaN(seconds) || seconds < 1 || seconds > SmallIntMax) {
              seconds = 0
            }
            setSuperSet('startSeconds', seconds)
          }}
        >
          <For each={Object.keys(StartType)}>
            {(typ) => (
              <option selected={typ === superSet.startType} value={typ}>
                <Trans key={typ.toLowerCase()} />
              </option>
            )}
          </For>
        </input>
      </div>

      <For each={superSet.subSets}>
        {(subSet, idx) => {
          return (
            <SetCard
              set={subSet}
              setNumber={idx() + 1}
              onSettingsClick={() => setSetSettingsOpener({ idx: idx() })}
            />
          )
        }}
      </For>

      <Show when={superSet.subSets?.length === 0}>
        <div class="m-4 rounded-lg bg-sky-200 p-4 text-xl font-semibold">
          <Trans key="no.sets.in.superset" />
        </div>
      </Show>
      <button
        class="float-right my-4 rounded-lg bg-sky-500 p-2 text-xl text-white shadow focus:outline-none focus:ring focus:ring-sky-300"
        onClick={() => setSetModalOpener({ set: newSubSet() })}
      >
        <Trans key="add.set" />
      </button>
      <div class="h-32 w-full"></div>
      <div class="fixed bottom-4 left-8 right-8 flex justify-between">
        <button
          class="rounded-lg bg-red-500 px-4 py-2 font-bold text-white hover:bg-sky-600 focus:outline-none focus:ring focus:ring-red-300"
          onClick={() => props.onClose()}
        >
          <Trans key="cancel" />
        </button>
        <button
          class="rounded-lg bg-green-500 px-4 py-2 font-bold text-white hover:bg-green-600 focus:outline-none focus:ring focus:ring-green-300"
          onClick={() => {
            if (!isSetValid()) return

            props.onSubmitSuperSet(unwrap(superSet))
          }}
        >
          <Trans key="add" />
        </button>
      </div>
    </div>
  )
}

export default SuperSetEditPage
