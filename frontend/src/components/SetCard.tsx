import { Component, For, Match, Show, Switch } from 'solid-js'
import { NewTrainingSet, StartType } from '../generated'
import settings from '../assets/settings.svg'
import { Trans } from '@mbarzda/solid-i18next'

type SetCardProps = {
  set: NewTrainingSet
  setNumber: number
  showSettings: boolean

  onSettingsClick?: (set: NewTrainingSet) => void
}

const SetCard: Component<SetCardProps> = (props) => {
  const setLayout = (setNum: number, set: NewTrainingSet) => {
    const repeat = set.repeat > 1 ? `${set.repeat} x ` : ''

    const startSeconds =
      set.startType !== StartType.None ? set.startSeconds! % 60 : 0
    const startMinutes =
      set.startType !== StartType.None
        ? (set.startSeconds! - startSeconds) / 60
        : 0
    const start = `${startMinutes !== 0 ? startMinutes + "'" : ''}${
      startSeconds !== 0 ? startSeconds + '"' : ''
    }`

    return (
      <div>
        <div class="rounded-t-lg bg-sky-300 p-2">
          <span class="inline-block w-10 text-xl">{setNum}.</span>
          <span class="inline-block w-10 text-xl">{repeat}</span>
          <span class="mr-4 text-xl">{set.distanceMeters}m</span>
          <Show when={set.startType !== 'None'}>
            <span class="text-xl">
              <Trans key={set.startType.toLowerCase()} />:{' '}
            </span>
            <span class="text-xl">{start}</span>
          </Show>
          <Show when={props.showSettings}>
            <img
              src={settings}
              width={32}
              height={32}
              class="float-right inline-block cursor-pointer"
              onClick={() => props.onSettingsClick?.(props.set)}
            />
          </Show>
        </div>
        <Show when={set.description}>
          <p class="whitespace-pre-wrap p-2 text-lg">{set.description}</p>
        </Show>
      </div>
    )
  }

  const superSetLayout = (setNum: number, superSet: NewTrainingSet) => {
    const repeat = superSet.repeat > 1 ? `${superSet.repeat} x ` : ''

    return (
      <div>
        <div class="rounded-t-lg bg-sky-300 p-2">
          <span class="inline-block w-10 text-xl">{setNum}.</span>
          <span class="inline-block w-10 text-xl">{repeat}</span>
          <Show when={props.showSettings}>
            <img
              src={settings}
              width={32}
              height={32}
              class="float-right inline-block cursor-pointer"
              onClick={() => props.onSettingsClick?.(props.set)}
            />
          </Show>
        </div>

        <For each={superSet.subSets}>
          {(subSet, idx) => {
            const repeatSubSet = subSet.repeat > 1 ? `${subSet.repeat} x ` : ''
            const startSeconds =
              subSet.startType !== StartType.None
                ? subSet.startSeconds! % 60
                : 0
            const startMinutes =
              subSet.startType !== StartType.None
                ? (subSet.startSeconds! - startSeconds) / 60
                : 0
            const start = `${startMinutes !== 0 ? startMinutes + "'" : ''}${
              startSeconds !== 0 ? startSeconds + '"' : ''
            }`

            return (
              <div class="ml-12 p-2">
                <span class="inline-block w-10 text-xl">
                  {(subSet.subSetOrder ?? 0) + 1}.
                </span>
                <span class="inline-block w-10 text-xl">{repeatSubSet}</span>
                <span class="mr-4 text-xl">{subSet.distanceMeters}m</span>
                <Show when={subSet.startType !== 'None'}>
                  <span class="text-xl">
                    <Trans key={subSet.startType.toLowerCase()} />:{' '}
                  </span>
                  <span class="text-xl">{start}</span>
                </Show>
                <p class="ml-10 whitespace-pre-wrap text-lg">
                  {subSet.description}
                </p>
                <Show when={idx() !== superSet.subSets!.length - 1}>
                  <hr class="mt-2 rounded-lg border-2 border-slate-500" />
                </Show>
              </div>
            )
          }}
        </For>
      </div>
    )
  }

  return (
    <div class="rounded-lg border border-solid border-slate-300 shadow md:mx-4 lg:mx-16 xl:mx-32">
      <Switch>
        <Match when={props.set.subSets && props.set.subSets?.length > 0}>
          {superSetLayout(props.setNumber, props.set)}
        </Match>
        <Match when={!props.set.subSets}>
          {setLayout(props.setNumber, props.set)}
        </Match>
      </Switch>
      <div class="rounded-b-lg border-t border-solid border-slate-300 p-2 text-lg">
        <span class="inline-block w-3/4">
          <Trans key="total.distance.set" />
        </span>
        <span class="inline-block w-1/4 text-end font-bold">
          {props.set.totalDistance}m
        </span>
      </div>
    </div>
  )
}

export default SetCard
