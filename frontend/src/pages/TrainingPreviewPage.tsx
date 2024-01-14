import { useTransContext } from '@mbarzda/solid-i18next'
import { Component, For, Show } from 'solid-js'
import DropdownMenu from '../components/common/DropdownMenu'
import ConfirmationModal from '../components/ConfirmationModal'
import { EquipmentIcons } from '../components/Equipment'
import { NewTraining, NewTrainingSet, StartType } from '../generated'

interface TrainingPreviewPageProps {
  training: NewTraining
  showOptions?: boolean
  options?: {
    onEdit: (setIdx: number) => void
    onDuplicate: (setIdx: number) => void
    onMoveUp: (setIdx: number) => void
    onMoveDown: (setIdx: number) => void
    onDelete: (setIdx: number) => void
  }
  showDeleteTraining?: boolean
  onDeleteTraining?: () => void
}

const TrainingPreviewPage: Component<TrainingPreviewPageProps> = (props) => {
  const [t] = useTransContext()

  const setCard = (set: NewTrainingSet) => {
    const setContent =
      set.repeat > 1
        ? `${set.repeat}x${set.distanceMeters}m`
        : `${set.distanceMeters}m`

    const startSeconds =
      set.startType !== StartType.None ? set.startSeconds! % 60 : 0
    const startMinutes =
      set.startType !== StartType.None
        ? (set.startSeconds! - startSeconds) / 60
        : 0
    const start = `${t(set.startType.toLowerCase())}: ${
      startMinutes !== 0 ? startMinutes + "'" : ''
    }${startSeconds !== 0 ? startSeconds + '"' : ''}`

    return (
      <div class="mx-auto block max-w-xl rounded-lg border border-gray-200 bg-white shadow">
        <div
          classList={{
            'rounded-lg':
              !set.description &&
              (!set.equipment || set.equipment.length === 0),
          }}
          class="w-full rounded-t-lg bg-sky-200 p-2 text-sky-900"
        >
          <h5 class="inline-block w-11/12 text-xl font-bold">
            <span class="pr-8">{setContent}</span>
            <Show when={set.startType !== StartType.None}>
              <span>{start}</span>
            </Show>
          </h5>

          <Show when={props.showOptions}>
            <DropdownMenu
              icon="fa-ellipsis"
              items={[
                {
                  text: t('edit'),
                  icon: 'fa-pen',
                  onClick: () => props.options!.onEdit(set.setOrder!),
                },
                {
                  text: t('duplicate'),
                  icon: 'fa-copy',
                  onClick: () => props.options!.onDuplicate(set.setOrder!),
                },
                {
                  text: t('move.up'),
                  icon: 'fa-arrow-up',
                  onClick: () => props.options!.onMoveUp(set.setOrder!),
                  disabled: set.setOrder === 0,
                },
                {
                  text: t('move.down'),
                  icon: 'fa-arrow-down',
                  onClick: () => props.options!.onMoveDown(set.setOrder!),
                  disabled: set.setOrder === props.training.sets.length - 1,
                },
              ]}
              finalItem={{
                text: t('delete'),
                icon: 'fa-trash',
                onClick: () => props.options!.onDelete(set.setOrder!),
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
        <div class="col-start-2 me-2 inline-block w-full rounded bg-sky-100 px-2.5 py-0.5 text-center text-xl font-medium text-sky-900">
          <span>{props.training.totalDistance / 1000}km</span>
        </div>
        <Show when={props.showDeleteTraining}>
          <div class="text-right">
            <ConfirmationModal
              icon="fa-trash"
              message={t('confirm.training.delete.message')}
              confirmLabel={t('confirm.delete.training')}
              cancelLabel={t('reject.delete.training')}
              onConfirm={props.onDeleteTraining!}
              onCancel={() => {}}
            />
          </div>
        </Show>
      </div>
      <div class="space-y-2">
        <For each={props.training.sets}>{setCard}</For>
      </div>
    </div>
  )
}

export default TrainingPreviewPage
