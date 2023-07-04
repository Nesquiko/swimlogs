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

interface TrainingPreviewProps {
  training?: Training | NewTraining
}

const TrainingPreview: Component<TrainingPreviewProps> = (props) => {
  /* const block = (block: Block | NewBlock) => { */
  /*   return ( */
  /*     <div class="mx-2"> */
  /*       <span class="text-2xl"> */
  /*         {block.num + 1}. {block.name} | <b>{block.totalDistance}m</b> */
  /*       </span> */
  /*       <div class="flex items-center justify-end"> */
  /*         <Show when={block.repeat > 1}> */
  /*           <span class="text-3xl font-bold">{block.repeat}x</span> */
  /*         </Show> */
  /*         <div class="w-11/12"> */
  /*           <For each={block.sets}>{set}</For> */
  /*         </div> */
  /*       </div> */
  /*     </div> */
  /*   ) */
  /* } */

  const set = (set: TrainingSet | NewTrainingSet) => {
    return (
      <div class="mx-8 my-2 flex flex-col space-x-4 rounded-lg border border-solid border-slate-300 bg-sky-50 p-2 shadow">
        <span class="text-lg">
          {set.repeat} x {set.distanceMeters}m
        </span>
        <Show when={set.startType !== StartType.None}>
          <span class="text-lg">
            {set.startType} {set.startSeconds}"
          </span>
        </Show>
        <span class="whitespace-pre-wrap text-lg">{set.description}</span>
      </div>
    )
  }

  const day = props.training?.start
    .toLocaleString('en', { weekday: 'long' })
    .toLowerCase()

  return (
    <div class=" m-2 space-y-2">
      <div class="mx-2">
        <p class="text-3xl">
          <Trans key={day!} />
        </p>
        <p class="text-3xl">{formatDate(props.training?.start)}</p>
        <p class="text-3xl">
          {props.training?.start.toString()} {props.training?.durationMin} min
        </p>
      </div>
      <For each={props.training?.sets}>{set}</For>
      <span class="text-2xl">
        <Trans key="distance.in.training" />{' '}
        <b>{props.training?.totalDistance}m</b>
      </span>
      <div class="h-32 w-full"></div>
    </div>
  )
}

export default TrainingPreview
