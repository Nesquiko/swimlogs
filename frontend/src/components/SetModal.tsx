import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { createEffect, For, JSX, on, onMount } from 'solid-js'
import { createStore, reconcile } from 'solid-js/store'
import { NewTrainingSet, StartType } from '../generated'
import { SmallIntMax } from '../lib/consts'

type SetModalProps = {
  opener: { set: NewTrainingSet; idx?: number }
  onSubmitSet: (set: NewTrainingSet, idx?: number) => void

  submitBtnLabelKey?: string
}

function SetModal(props: SetModalProps): JSX.Element {
  let dialog: HTMLDialogElement

  createEffect(
    on(
      () => props.opener,
      () => setTrainingSet(reconcile(props.opener.set)),
      { defer: true }
    )
  )

  const [trainingSet, setTrainingSet] = createStore<NewTrainingSet>(
    props.opener.set
  )

  const [t] = useTransContext()
  const submitLabel = props.submitBtnLabelKey
    ? t(props.submitBtnLabelKey)
    : t('add', 'Add')

  createEffect(
    on(
      () => props.opener,
      () => {
        dialog.inert = true // disables auto focus when dialog is opened
        dialog.showModal()
        dialog.inert = false
      },
      { defer: true }
    )
  )

  onMount(() => {
    dialog.addEventListener('click', (e) => {
      const dialogDimensions = dialog.getBoundingClientRect()
      if (
        e.clientX < dialogDimensions.left ||
        e.clientX > dialogDimensions.right ||
        e.clientY < dialogDimensions.top ||
        e.clientY > dialogDimensions.bottom
      ) {
        dialog.close()
      }
    })
  })

  const isSetValid = (): boolean => {
    if (trainingSet.repeat < 1) {
      return false
    }

    if (
      trainingSet.distanceMeters === undefined ||
      trainingSet.distanceMeters < 1
    ) {
      return false
    }

    if (
      trainingSet.startType !== StartType.None &&
      trainingSet.startSeconds !== undefined &&
      trainingSet.startSeconds < 1
    ) {
      return false
    }

    return true
  }

  return (
    <dialog ref={dialog!} class="h-screen w-11/12 rounded-lg">
      <p class="text-center text-2xl">
        <Trans key="add.new.set" />
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
            'border-red-500 text-red-500': trainingSet.repeat < 1,
            'border-slate-300': trainingSet.repeat >= 1
          }}
          class="w-24 rounded-md border p-2 text-center text-lg focus:border-blue-500 focus:outline-none focus:ring"
          value={trainingSet.repeat}
          onChange={(e) => {
            let repeat = parseInt(e.target.value)
            if (Number.isNaN(repeat) || repeat < 1 || repeat > SmallIntMax) {
              repeat = 0
            }
            setTrainingSet('repeat', repeat)
            setTrainingSet(
              'totalDistance',
              repeat * (trainingSet.distanceMeters ?? 0)
            )
          }}
        />
      </div>
      <div class="my-2 flex items-center justify-between">
        <label class="text-xl" for="distance">
          <Trans key="distance" />
        </label>
        <input
          id="distance"
          type="number"
          placeholder="400"
          classList={{
            'border-red-500 text-red-500':
              trainingSet.distanceMeters !== undefined &&
              trainingSet.distanceMeters < 1,
            'border-slate-300':
              trainingSet.distanceMeters === undefined ||
              trainingSet.distanceMeters >= 1
          }}
          class="w-24 rounded-md border p-2 text-center text-lg focus:border-blue-500 focus:outline-none focus:ring"
          value={trainingSet.distanceMeters ?? '0'}
          onChange={(e) => {
            let dist = parseInt(e.target.value)
            if (Number.isNaN(dist) || dist < 1 || dist > SmallIntMax) {
              dist = 0
            }
            setTrainingSet('distanceMeters', dist)
            setTrainingSet('totalDistance', trainingSet.repeat * dist)
          }}
        />
      </div>
      <div class="my-2 flex items-center justify-between">
        <label class="text-xl" for="start">
          <Trans key="start" />
        </label>
        <select
          id="start"
          class="w-32 rounded-md border border-solid border-slate-300 bg-white p-2 text-start text-lg focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            setTrainingSet('startType', e.target.value as StartType)
          }}
        >
          <For each={Object.keys(StartType)}>
            {(typ) => (
              <option selected={typ === trainingSet.startType} value={typ}>
                <Trans key={typ.toLowerCase()} />
              </option>
            )}
          </For>
        </select>
      </div>
      <div
        classList={{
          visible: trainingSet.startType !== StartType.None,
          invisible: trainingSet.startType === StartType.None
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
              trainingSet.startSeconds !== undefined &&
              trainingSet.startSeconds < 1,
            'border-slate-300':
              trainingSet.startSeconds === undefined ||
              trainingSet.startSeconds >= 1
          }}
          class="w-24 rounded-md border border-solid border-slate-300 bg-white p-2 text-center text-lg focus:border-sky-500 focus:outline-none focus:ring"
          value={trainingSet.startSeconds ?? ''}
          onChange={(e) => {
            const val = e.target.value
            let seconds = parseInt(val)
            if (Number.isNaN(seconds) || seconds < 1 || seconds > SmallIntMax) {
              seconds = 0
            }
            setTrainingSet('startSeconds', seconds)
          }}
        >
          <For each={Object.keys(StartType)}>
            {(typ) => (
              <option selected={typ === trainingSet.startType} value={typ}>
                <Trans key={typ.toLowerCase()} />
              </option>
            )}
          </For>
        </input>
      </div>
      <label class="text-xl" for="description">
        <Trans key="description" />
      </label>
      <textarea
        id="description"
        placeholder={t(
          'set.description.placeholder',
          'Freestyle, Breastroke, etc.'
        )}
        maxlength="255"
        class="w-full rounded-lg border p-2 text-lg focus:border-sky-500 focus:outline-none focus:ring"
        value={trainingSet.description ?? ''}
        onChange={(e) => {
          let desc: string | undefined = e.target.value
          if (desc === '') {
            desc = undefined
          }
          setTrainingSet('description', desc)
        }}
      />
      <div class="absolute bottom-4 left-8 right-8 flex justify-between">
        <button
          class="rounded-lg bg-red-500 px-4 py-2 font-bold text-white hover:bg-sky-600 focus:outline-none focus:ring focus:ring-red-300"
          onClick={() => dialog.close()}
        >
          <Trans key="cancel" />
        </button>
        <button
          class="rounded-lg bg-green-500 px-4 py-2 font-bold text-white hover:bg-green-600 focus:outline-none focus:ring focus:ring-green-300"
          onClick={() => {
            if (!isSetValid()) return
            // this modal only creates sets without subsets, so we can just shallow copy the object
            props.onSubmitSet(Object.assign({}, trainingSet), props.opener.idx)
            dialog.close()
          }}
        >
          {submitLabel}
        </button>
      </div>
    </dialog>
  )
}

export default SetModal
