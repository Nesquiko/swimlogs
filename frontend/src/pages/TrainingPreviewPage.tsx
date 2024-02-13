import { useTransContext } from '@mbarzda/solid-i18next'
import { type Component, For, Show, type JSX } from 'solid-js'
import DropdownMenu from '../components/common/DropdownMenu'
import { EquipmentIcons } from '../components/Equipment'
import {
  EquipmentEnum,
  NewTraining,
  NewTrainingSet,
  StartTypeEnum,
  Training,
} from 'swimlogs-api'
import SetCard, { Option, SkeletonSetCard } from '../components/SetCard'

interface TrainingPreviewPageProps {
  training: NewTraining | Training

  setOptions?: Option[]

  rightHeaderComponent?: Component
  leftHeaderComponent?: Component
}

const TrainingPreviewPage: Component<TrainingPreviewPageProps> = (props) => {
  const [t] = useTransContext()

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
        <For each={props.training.sets}>
          {(set) => <SetCard set={set} setOptions={props.setOptions} />}
        </For>
      </div>
    </div>
  )
}

const SkeletonTrainingPreviewPage: Component = () => {
  return (
    <div class="space-y-4 px-4 animate-pulse">
      <div class="grid grid-cols-3 items-center">
        <div class="col-start-2 h-8 me-2 inline-block w-full rounded bg-sky-100 px-2.5 py-0.5 text-center text-sky-900"></div>
      </div>
      <div class="space-y-2">
        <For each={Array(8)}>{(_) => <SkeletonSetCard />}</For>
      </div>
    </div>
  )
}

export default TrainingPreviewPage
export { SkeletonTrainingPreviewPage }
