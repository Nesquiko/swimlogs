import { useTransContext } from '@mbarzda/solid-i18next'
import { TFunction } from 'i18next'
import { Component, For } from 'solid-js'
import { Session } from '../generated'
import Pagination from './Pagination'

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
      <Pagination
        onPrevPage={() => props.onPrevPage?.()}
        prevDisabled={props.sessionPage === 0}
        onNextPage={() => props.onNextPage?.()}
        nextDisabled={props.isLastPage}
      />
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
