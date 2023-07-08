import { Component, Show } from 'solid-js'
import { NewTrainingSet } from '../generated'
import settings from '../assets/settings.svg'
import { Trans } from '@mbarzda/solid-i18next'

type SetCardProps = {
  set: NewTrainingSet
  onSettingsClick: (set: NewTrainingSet) => void
}

const SetCard: Component<SetCardProps> = (props) => {
  return (
    <div class="my-2 rounded-lg border border-solid border-slate-300 bg-sky-100 p-2 shadow">
      <span class="inline-block w-10 text-xl">
        {(props.set.subSetOrder ?? 0) + 1}.
      </span>
      <Show when={props.set.repeat > 1}>
        <span class="text-xl">{props.set.repeat} x </span>
      </Show>
      <span class="mr-4 text-xl">{props.set.distanceMeters}m</span>
      <Show when={props.set.startType !== 'None'}>
        <span class="text-xl">
          <Trans key={props.set.startType.toLowerCase()} />:{' '}
        </span>
        <span class="text-xl">{props.set.startSeconds}"</span>
      </Show>
      <img
        src={settings}
        width={32}
        height={32}
        class="float-right cursor-pointer"
        onClick={() => props.onSettingsClick(props.set)}
      />
      <p class="mx-6 whitespace-pre-wrap text-xl">{props.set.description}</p>
    </div>
  )
}

export default SetCard
