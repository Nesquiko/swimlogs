import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { Component, For, onMount, Show } from 'solid-js'
import { InvalidTrainingSet, NewTrainingSet, StartType } from '../generated'
import { useCreateTraining } from './context/CreateTrainingContextProvider'
import settings from '../assets/settings.svg'
import { SmallIntMax } from '../lib/consts'
import { isInvalidSetEmpty } from '../lib/validation'

interface SetFormProps {
  set: NewTrainingSet
  setIdx: number
  invalidSet: InvalidTrainingSet
  onDuplicate: () => void
  onDelete: () => void
}

// TODO instead of delete, duplicate buttons, create a menu with these
export const SetForm: Component<SetFormProps> = (props) => {
  const [{ setTraining }, { setInvalidTraining }, , , ,] = useCreateTraining()
  const [t] = useTransContext()

  return (
    <div
      classList={{
        'bg-sky-50': isInvalidSetEmpty(props.invalidSet),
        'bg-red-50': !isInvalidSetEmpty(props.invalidSet)
      }}
      class="mx-auto my-2 rounded-lg border border-solid border-slate-300 px-2 shadow"
    >
      <div class="my-2">
        <input
          type="number"
          placeholder="1"
          classList={{
            'border-red-500 text-red-500':
              props.invalidSet.repeat !== undefined,
            'border-slate-300': props.invalidSet.repeat === undefined
          }}
          class="w-1/6 rounded-md border p-2 text-center focus:border-blue-500 focus:outline-none focus:ring"
          value={props.set.repeat}
          onChange={(e) => {
            let repeat = parseInt(e.target.value)
            if (Number.isNaN(repeat) || repeat < 1 || repeat > SmallIntMax) {
              repeat = 0
            }
            setTraining('sets', props.setIdx, 'repeat', repeat)
            setTraining(
              'sets',
              props.setIdx,
              'totalDistance',
              repeat * (props.set.distanceMeters ?? 0)
            )
            setInvalidTraining(
              'invalidSets',
              props.setIdx,
              'repeat',
              repeat > 0
                ? undefined
                : 'Repeat must be a number between 1 and 32767'
            )
          }}
        />
        <span class="mx-4 text-lg">x</span>
        <input
          type="number"
          placeholder="400"
          classList={{
            'border-red-500 text-red-500':
              props.invalidSet.distanceMeters !== undefined,
            'border-slate-300': props.invalidSet.distanceMeters === undefined
          }}
          class="w-1/3 rounded-md border p-2 text-center focus:border-blue-500 focus:outline-none focus:ring"
          value={props.set.distanceMeters}
          onChange={(e) => {
            let distance = parseInt(e.target.value)
            if (
              Number.isNaN(distance) ||
              distance < 1 ||
              distance > SmallIntMax
            ) {
              distance = 0
            }
            setTraining('sets', props.setIdx, 'distanceMeters', distance)
            setTraining(
              'sets',
              props.setIdx,
              'totalDistance',
              props.set.repeat * distance
            )
            setInvalidTraining(
              'invalidSets',
              props.setIdx,
              'distanceMeters',
              distance > 0
                ? undefined
                : 'Distance must be a number between 1 and 32767'
            )
          }}
        />
        <Dropdown
          items={[
            {
              label: t('duplicate', 'Duplicate'),
              action: () => props.onDuplicate()
            },
            { label: t('delete', 'Delete'), action: () => props.onDelete() }
          ]}
        />
      </div>
      <textarea
        placeholder={t(
          'set.description.placeholder',
          'Freestyle, Breastroke, etc.'
        )}
        maxlength="255"
        class="w-full rounded-md border p-2 focus:border-sky-500 focus:outline-none focus:ring"
        value={props.set.description ?? ''}
        onChange={(e) => {
          let desc: string | undefined = e.target.value
          if (desc === '') {
            desc = undefined
          }
          setTraining('sets', props.setIdx, 'description', desc)
        }}
      />
      <div class="flex items-center space-x-4">
        <label for="start">
          <Trans key="start" />
        </label>
        <select
          id="start"
          class="my-2 rounded-md border border-solid border-slate-300 bg-white p-2 focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            setTraining(
              'sets',
              props.setIdx,
              'startType',
              e.target.value as StartType
            )
          }}
        >
          <For each={Object.keys(StartType)}>
            {(typ) => (
              <option selected={typ === props.set.startType} value={typ}>
                <Trans key={typ.toLowerCase()} />
              </option>
            )}
          </For>
        </select>
        <Show when={props.set.startType !== StartType.None}>
          <label for="seconds">
            <Trans key="seconds" />
          </label>
          <input
            id="seconds"
            type="number"
            placeholder="20"
            classList={{
              'border-red-500 text-red-500':
                props.invalidSet.startSeconds !== undefined,
              'border-slate-300': props.invalidSet.startSeconds === undefined
            }}
            class="w-1/4 rounded-md border border-solid px-4 py-2 text-center focus:border-sky-500 focus:outline-none focus:ring"
            value={props.set.startSeconds}
            onChange={(e) => {
              const val = e.target.value
              let seconds = parseInt(val)
              if (
                Number.isNaN(seconds) ||
                seconds < 1 ||
                seconds > SmallIntMax
              ) {
                seconds = 0
              }
              setTraining('sets', props.setIdx, 'startSeconds', seconds)
              setInvalidTraining(
                'invalidSets',
                props.setIdx,
                'startSeconds',
                seconds > 0
                  ? undefined
                  : 'Seconds must be a number between 1 and 32767'
              )
            }}
          />
        </Show>
      </div>
    </div>
  )
}

type DropdownItem = {
  label: string
  action: () => void
}

type DropdownProps = {
  items: DropdownItem[]
}

const Dropdown: Component<DropdownProps> = (props) => {
  let dialog: HTMLDialogElement

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

  return (
    <div class="float-right">
      <img
        src={settings}
        width={32}
        height={32}
        class="cursor-pointer"
        onClick={() => {
          dialog.inert = true // disables auto focus when dialog is opened
          dialog.showModal()
          dialog.inert = false
        }}
      />
      <dialog ref={dialog!} class="w-44 rounded-lg">
        <ul class="text-md text-black">
          <For each={props.items}>
            {(item) => {
              return (
                <li>
                  <p
                    class="block cursor-pointer p-2 hover:bg-slate-200"
                    onClick={() => {
                      item.action()
                      dialog.close()
                    }}
                  >
                    {item.label}
                  </p>
                </li>
              )
            }}
          </For>
        </ul>
      </dialog>
    </div>
  )
}

export default SetForm
