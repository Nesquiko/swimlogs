import { useTransContext } from '@mbarzda/solid-i18next'
import { type Component, For, Show, type JSX } from 'solid-js'
import DropdownMenu from '../components/common/DropdownMenu'
import { EquipmentIcons } from '../components/Equipment'
import {
  NewTraining,
  NewTrainingSet,
  StartTypeEnum,
  Training,
} from 'swimlogs-api'

interface Option {
  text: string
  icon: string
  onClick: (setIdx: number) => void
  disabled?: boolean
  disabledFunc?: (setIdx: number) => boolean
}

interface TrainingPreviewPageProps {
  training: NewTraining | Training

  setOptions?: Option[]

  rightHeaderComponent?: Component
  leftHeaderComponent?: Component
}

const TrainingPreviewPage: Component<TrainingPreviewPageProps> = (props) => {
  const [t] = useTransContext()

  const setCard = (set: NewTrainingSet): JSX.Element => {
    const setContent =
      set.repeat > 1
        ? `${set.repeat}x${set.distanceMeters}m`
        : `${set.distanceMeters}m`

    const start = (): string => {
      if (set.startSeconds === undefined) {
        return ''
      }

      const startSeconds =
        set.startType !== StartTypeEnum.None ? set.startSeconds % 60 : 0
      const startMinutes =
        set.startType !== StartTypeEnum.None
          ? (set.startSeconds - startSeconds) / 60
          : 0
      return `${t(set.startType.toLowerCase())}: ${
        startMinutes !== 0 ? startMinutes + "'" : ''
      }${startSeconds !== 0 ? startSeconds + '"' : ''}`
    }

    return (
      <div class="mx-auto block max-w-xl rounded-lg border border-gray-200 bg-white shadow">
        <div
          classList={{
            'rounded-lg':
              set.description === undefined &&
              (set.equipment === undefined || set.equipment.length === 0),
          }}
          class="w-full rounded-t-lg bg-sky-200 p-2 text-sky-900"
        >
          <h5 class="inline-block w-11/12 text-xl font-bold">
            <span class="pr-8">{setContent}</span>
            <Show when={set.startType !== StartTypeEnum.None}>
              <span>{start()}</span>
            </Show>
          </h5>

          <Show
            when={props.setOptions !== undefined && props.setOptions.length > 0}
          >
            <DropdownMenu
              icon="fa-ellipsis"
              items={(props.setOptions ?? []).slice(0, -1).map((option) => ({
                icon: option.icon,
                text: option.text,
                onClick: () => option.onClick(set.setOrder!),
                disabled:
                  option.disabled ||
                  option.disabledFunc?.(set.setOrder!) ||
                  false,
              }))}
              finalItem={{
                icon: props.setOptions![props.setOptions!.length - 1].icon,
                text: props.setOptions![props.setOptions!.length - 1].text,
                onClick: () =>
                  props.setOptions![props.setOptions!.length - 1].onClick(
                    set.setOrder!
                  ),
                disabled:
                  props.setOptions![props.setOptions!.length - 1].disabled ||
                  props.setOptions![
                    props.setOptions!.length - 1
                  ].disabledFunc?.(set.setOrder!) ||
                  false,
              }}
            />
          </Show>
        </div>
        <Show when={set.description}>
          <p class="whitespace-pre-wrap p-2 text-gray-500">{set.description}</p>
        </Show>
        <Show when={set.equipment && set.equipment.length > 0}>
          <div class="text-center">
            <For each={set.equipment}>
              {(equipment) => (
                <img
                  class="inline-block"
                  width={64}
                  height={64}
                  src={EquipmentIcons.get(equipment)}
                />
              )}
            </For>
          </div>
        </Show>
      </div>
    )
  }

  return (
    <div class="space-y-4 px-4">
      <div class="grid grid-cols-3 items-center">
        <Show when={props.leftHeaderComponent}>
          {props.leftHeaderComponent}
        </Show>
        <div class="col-start-2 me-2 inline-block w-full rounded bg-sky-100 px-2.5 py-0.5 text-center text-xl font-medium text-sky-900">
          <span>{props.training.totalDistance / 1000}km</span>
        </div>
        <Show when={props.rightHeaderComponent}>
          {props.rightHeaderComponent}
        </Show>
      </div>
      <div class="space-y-2">
        <For each={props.training.sets}>{setCard}</For>
      </div>
    </div>
  )
}

export default TrainingPreviewPage
