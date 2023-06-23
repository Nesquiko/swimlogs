import { Component, For, Show } from 'solid-js'
import {
  Block,
  NewBlock,
  NewTraining,
  NewTrainingSet,
  StartingRuleType,
  Training,
  TrainingSet
} from '../generated'

interface TrainingPreviewProps {
  training?: Training | NewTraining
}

const TrainingPreview: Component<TrainingPreviewProps> = (props) => {
  const block = (block: Block | NewBlock) => {
    return (
      <div class="flex flex-col p-2">
        <span class="text-2xl">
          Block {block.num + 1}. {block.name} | <b>{block.totalDistance}m</b>
        </span>
        <span class="text-xl">{block.repeat}x</span>
        <For each={block.sets}>{set}</For>
      </div>
    )
  }

  const set = (set: TrainingSet | NewTrainingSet) => {
    return (
      <div class="mx-8 my-2 flex space-x-4 rounded-lg border border-solid border-slate-300 bg-sky-50 p-2 shadow">
        <span class="text-lg">Set {set.num + 1}.</span>
        <div class="flex flex-col">
          <span class="text-lg">
            {set.repeat} x {set.distance}m
          </span>
          <span class="whitespace-pre-wrap text-lg">{set.what}</span>
          <Show when={set.startingRule.type !== StartingRuleType.None}>
            <span class="text-lg">
              Start {set.startingRule.type} {set.startingRule.seconds}"
            </span>
          </Show>
        </div>
      </div>
    )
  }

  return (
    <div class="m-2">
      <h1 class="text-2xl">
        Date{' '}
        <b>
          {props.training?.date.toLocaleDateString('sk-SK').replaceAll(' ', '')}
        </b>
      </h1>
      <h1 class="text-2xl">
        Day{' '}
        <b>{props.training?.date.toLocaleString('en', { weekday: 'long' })}</b>
      </h1>
      <h1 class="text-2xl">
        Start time <b>{props.training?.startTime}</b>
      </h1>
      <For each={props.training?.blocks}>{block}</For>
      <span class="text-2xl">
        Distance in training <b>{props.training?.totalDistance}m</b>
      </span>
      <div class="h-32 w-full"></div>
    </div>
  )
}

export default TrainingPreview
