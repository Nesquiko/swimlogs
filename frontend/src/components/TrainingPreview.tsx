import { Trans } from '@mbarzda/solid-i18next'
import { Component, For, Show } from 'solid-js'
import {
  NewTraining,
  NewTrainingSet,
  StartType,
  Training,
  TrainingSet
} from '../generated'
import { formatDate } from '../lib/datetime'
import SetCard from './SetCard'

interface TrainingPreviewProps {
  training?: Training | NewTraining
}

const TrainingPreview: Component<TrainingPreviewProps> = (props) => {
  const day = props.training?.start
    .toLocaleString('en', { weekday: 'long' })
    .toLowerCase()

  const time = props.training?.start.toLocaleTimeString('sk-SK', {
    hour: '2-digit',
    minute: '2-digit'
  })

  return (
    <div class="m-2 space-y-4">
      <div class="mx-2">
        <p class="text-2xl">
          <Trans key={day!} /> - {formatDate(props.training?.start)}
        </p>
        <p class="text-2xl">
          <Trans key="starttime" /> {time}
        </p>
        <p class="text-2xl">
          <Trans key="duration" /> {props.training?.durationMin} min
        </p>
        <p class="text-2xl">
          <Trans key="distance.in.training" />{' '}
          <b>{props.training?.totalDistance}m</b>
        </p>
      </div>
      <For each={props.training?.sets}>
        {(set, idx) => (
          <SetCard set={set} setNumber={idx() + 1} showSettings={false} />
        )}
      </For>
      <div class="h-32 w-full"></div>
    </div>
  )
}

export default TrainingPreview
