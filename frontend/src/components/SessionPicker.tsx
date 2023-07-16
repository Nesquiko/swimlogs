import { useTransContext } from '@mbarzda/solid-i18next'
import { TFunction } from 'i18next'
import { Component, For } from 'solid-js'
import { Session } from '../generated'

interface SessionPickerProps {
  sessions?: Session[]
  selectedSession: Session | 'not-selected' | undefined
  sessionPage: number
  isLastPage: boolean
  onNextPage?: () => void
  onPrevPage?: () => void
  onSelect: (s: Session) => void
}

const SessionPicker: Component<SessionPickerProps> = (props) => {
  const [t] = useTransContext()

  return (
    <div class="h-full">
      <div class="h-3/4 space-y-2 rounded-lg px-4 md:px-8 lg:px-12 xl:px-16">
        <For each={props.sessions}>
          {(s) => (
            <SessionPickerItem
              s={s}
              t={t}
              selected={s === props.selectedSession}
              onClick={() => props.onSelect(s)}
            />
          )}
        </For>
      </div>
      <div class="space-x-8 p-2 text-center">
        <button
          classList={{
            'text-slate-300 pointer-events-none': props.sessionPage === 0
          }}
          disabled={props.sessionPage === 0}
          class="inline-flex cursor-pointer text-black"
          onClick={() => props.onPrevPage?.()}
        >
          <svg
            aria-hidden="true"
            class="h-8 w-8"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
              clip-rule="evenodd"
            ></path>
          </svg>
        </button>
        <button
          classList={{ 'text-slate-300 pointer-events-none': props.isLastPage }}
          class="inline-flex cursor-pointer text-black"
          disabled={props.isLastPage}
          onClick={() => props.onNextPage?.()}
        >
          <svg
            aria-hidden="true"
            class="h-8 w-8"
            fill="currentColor"
            viewBox="0 0 20 20"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill-rule="evenodd"
              d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
              clip-rule="evenodd"
            ></path>
          </svg>
        </button>
      </div>
    </div>
  )
}

interface SessionPickerItemProps {
  s: Session
  t: TFunction
  selected: boolean
  onClick: () => void
}

const SessionPickerItem: Component<SessionPickerItemProps> = (props) => {
  const localizedDay = props.t(props.s.day.toLowerCase(), props.s.day)

  return (
    <div
      classList={{ 'bg-sky-500/20 border-slate-400': props.selected }}
      class="flex w-full cursor-pointer rounded-lg border border-solid p-2 text-center shadow"
      onClick={props.onClick}
    >
      <span class="w-1/3 text-lg">{localizedDay}</span>
      <span class="w-1/3 text-lg">{props.s.startTime}</span>
      <span class="w-1/3 text-lg">{props.s.durationMin} min</span>
    </div>
  )
}

export default SessionPicker
