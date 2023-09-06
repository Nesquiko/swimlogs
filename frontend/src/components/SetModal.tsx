import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { createEffect, createSignal, For, JSX, on, onMount } from 'solid-js'
import { createStore, produce, reconcile } from 'solid-js/store'
import { Equipment, NewTrainingSet, StartType } from '../generated'
import { SmallIntMax } from '../lib/consts'
import Counter from './Counter'
import { EquipmentIcons } from './Equipment'

type SetModalProps = {
  opener: { set: NewTrainingSet; idx?: number }
  onSubmitSet: (set: NewTrainingSet, idx?: number) => void
  title: string
  submitBtnLabelKey?: string
}

function SetModal(props: SetModalProps): JSX.Element {
  let dialog: HTMLDialogElement
  const [startTimeSeconds, setStartTimeSeconds] = createSignal(0)
  const [startTimeSecondsErr, setStartTimeSecondsErr] = createSignal(false)

  const [startTimeMinutes, setStartTimeMinutes] = createSignal(0)
  const [startTimeMinutesErr, setStartTimeMinutesErr] = createSignal(false)

  createEffect(
    on(
      () => props.opener,
      () => {
        setStartTimeMinutes(0)
        setStartTimeSeconds(0)
        setStartTimeMinutesErr(false)
        setStartTimeSecondsErr(false)

        if (props.opener.set.startType !== StartType.None) {
          const seconds = props.opener.set.startSeconds! % 60
          const minutes = (props.opener.set.startSeconds! - seconds) / 60
          setStartTimeMinutes(minutes)
          setStartTimeSeconds(seconds)
        }
        setTrainingSet(reconcile(props.opener.set))
      },
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

  const submit = () => {
    if (!isSetValid()) return

    // this modal only creates sets without subsets, so we can just shallow copy the object
    const ts = Object.assign({}, trainingSet)
    if (ts.startType !== StartType.None) {
      ts.startSeconds = startTimeSeconds() + startTimeMinutes() * 60
    } else {
      ts.startSeconds = undefined
    }
    props.onSubmitSet(ts, props.opener.idx)
    dialog.close()
  }

  const isSetValid = (): boolean => {
    let isValid = true
    if (trainingSet.repeat < 1) {
      isValid = false
    }

    if (
      trainingSet.distanceMeters === undefined ||
      trainingSet.distanceMeters < 1
    ) {
      isValid = false
    }

    if (trainingSet.startType !== StartType.None) {
      if (startTimeSeconds() < 0 || startTimeSeconds() > 59) {
        setStartTimeSecondsErr(true)
        isValid = false
      }

      if (startTimeMinutes() < 0 || startTimeMinutes() > 59) {
        setStartTimeMinutesErr(true)
        isValid = false
      }

      if (startTimeMinutes() * 60 + startTimeSeconds() === 0) {
        setStartTimeMinutesErr(true)
        setStartTimeSecondsErr(true)
        isValid = false
      }
    }

    return isValid
  }

  const equipmentButton = (equipment: Equipment) => {
    let imgSrc = EquipmentIcons.get(equipment)
    return (
      <button
        classList={{
          'bg-sky-300': trainingSet.equipment?.includes(equipment)
        }}
        class="rounded-lg border border-slate-300 p-2"
        onClick={() => {
          if (trainingSet.equipment === undefined) {
            setTrainingSet('equipment', new Array())
          }

          if (trainingSet.equipment?.includes(equipment)) {
            setTrainingSet('equipment', (e) =>
              e?.filter((e) => e !== equipment)
            )
            return
          }

          setTrainingSet(
            'equipment',
            produce((e) => {
              e?.push(equipment)
            })
          )
        }}
      >
        <img width={48} height={48} src={imgSrc} />
      </button>
    )
  }

  const setDistance = (distance: number) => {
    setTrainingSet('distanceMeters', distance)
    setTrainingSet('totalDistance', trainingSet.repeat * distance)
  }

  return (
    <dialog
      ref={dialog!}
      class="h-screen w-11/12 rounded-lg md:w-5/6 lg:w-2/3 xl:w-1/3"
    >
      <p class="text-center text-2xl">{props.title}</p>
      <hr class="my-4 rounded-lg border-2 border-slate-500" />
      <div class="flex items-center justify-between">
        <label class="text-xl">
          <Trans key="repeat" />
        </label>
        <Counter
          initial={trainingSet.repeat}
          min={1}
          max={SmallIntMax}
          onChange={(n) => {
            setTrainingSet('repeat', n)
            setTrainingSet(
              'totalDistance',
              n * (trainingSet.distanceMeters ?? 0)
            )
          }}
        />
      </div>
      <p class="my-4 text-xl">
        <Trans key="distance" />
      </p>
      <div class="my-4 flex items-center justify-between">
        <button
          classList={{
            'bg-sky-400 text-white': trainingSet.distanceMeters === 25
          }}
          class="w-16 rounded-lg border border-slate-300 p-2 text-center text-lg"
          onClick={() => setDistance(25)}
        >
          25
        </button>
        <button
          classList={{
            'bg-sky-400 text-white': trainingSet.distanceMeters === 50
          }}
          class="w-16 rounded-lg border border-slate-300 p-2 text-center text-lg"
          onClick={() => setDistance(50)}
        >
          50
        </button>
        <button
          classList={{
            'bg-sky-400 text-white': trainingSet.distanceMeters === 100
          }}
          class="w-16 rounded-lg border border-slate-300 p-2 text-center text-lg"
          onClick={() => setDistance(100)}
        >
          100
        </button>
        <input
          id="distance"
          type="number"
          placeholder={t('different.distance', 'Custom distance')}
          classList={{
            'border-red-500 text-red-500':
              trainingSet.distanceMeters !== undefined &&
              trainingSet.distanceMeters < 1,
            'border-slate-300':
              trainingSet.distanceMeters === undefined ||
              trainingSet.distanceMeters >= 1
          }}
          class="w-36 rounded-lg border p-2 text-center text-lg focus:border-blue-500 focus:outline-none focus:ring"
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
      <div class="my-4 flex items-center justify-between">
        <label class="text-xl" for="start">
          <Trans key="start" />
        </label>
        <select
          id="start"
          class="w-32 rounded-lg border border-solid border-slate-300 bg-white p-2 text-start text-lg focus:border-sky-500 focus:outline-none focus:ring"
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
        class="my-4 flex justify-around gap-4 md:justify-end"
      >
        <Counter
          label={t('minutes')}
          initial={startTimeMinutes()}
          onChange={(m) => setStartTimeMinutes(m)}
          min={0}
          max={59}
          error={startTimeMinutesErr()}
        />
        <Counter
          label={t('seconds')}
          initial={startTimeSeconds()}
          onChange={(s) => setStartTimeSeconds(s)}
          min={0}
          max={59}
          error={startTimeSecondsErr()}
        />
      </div>
      <div class="my-4 ">
        <label class="text-xl">
          <Trans key="equipment" />
        </label>
        <div class="flex items-center justify-between">
          {equipmentButton(Equipment.Fins)}
          {equipmentButton(Equipment.Monofin)}
          {equipmentButton(Equipment.Snorkel)}
          {equipmentButton(Equipment.Paddles)}
          {equipmentButton(Equipment.Board)}
        </div>
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
          class="rounded-lg bg-red-500 px-4 py-2 font-bold text-white hover:bg-red-600 focus:outline-none focus:ring focus:ring-red-300"
          onClick={() => dialog.close()}
        >
          <Trans key="cancel" />
        </button>
        <button
          class="rounded-lg bg-sky-500 px-4 py-2 font-bold text-white hover:bg-sky-600 focus:outline-none focus:ring focus:ring-sky-300"
          onClick={submit}
        >
          {submitLabel}
        </button>
      </div>
    </dialog>
  )
}

export default SetModal
