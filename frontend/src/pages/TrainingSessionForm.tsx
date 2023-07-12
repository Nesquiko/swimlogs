import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { A } from '@solidjs/router'
import {
  Component,
  createEffect,
  createSignal,
  For,
  Match,
  Show,
  Switch
} from 'solid-js'
import SessionPicker from '../components/SessionPicker'
import { useSessionsContext } from '../components/SessionsContextProvider'
import { useShownComponent } from '../components/ShownComponentContextProvider'
import { openToast, ToastType } from '../components/Toast'
import { useStateContext } from '../components/TrainingStateContext'
import { Day, Session } from '../generated'
import {
  NullDateTime,
  NullDay,
  NullStartTime,
  SmallIntMax,
  StartTimeHours,
  StartTimeMinutes
} from '../lib/consts'
import { formatDate } from '../lib/datetime'

export const TrainingSessionForm: Component = () => {
  const [t] = useTransContext()
  const [{ training, setTraining }, state] = useStateContext()
  const [, setCurrentComponent] = useShownComponent()
  const [
    sessions,
    fetchNextSessionPage,
    fetchPrevSessionPage,
    sessionsPage,
    isLastPage,
    serverError
  ] = useSessionsContext()

  createEffect(() => {
    if (serverError() !== undefined) {
      setPickSession(false)
      openToast(
        t(
          'session.fetch.error',
          'There was an error during loading sessions, please set session manually'
        ),
        ToastType.ERROR,
        4000
      )
    }
  })

  const [pickSession, setPickSession] = createSignal(true)
  const [day, setDay] = state.day
  const [durationMin, setDurationMin] = state.durationMin
  const [startTime, setStartTime] = state.startTime
  const [dates] = state.dates
  const [selectedDate, setSelectedDate] = state.selectedDate
  const [session, setSelectedSession] = state.session

  const sumbit = () => {
    let isValid = true

    if (pickSession()) {
      isValid = validateWhenPickingSession()
    } else {
      isValid = validateWhenSettingSession()
    }

    if (selectedDate() === NullDateTime) {
      setSelectedDate(undefined)
      isValid = false
    }

    if (!isValid) {
      return
    }

    const s = session() as Session
    setTraining('durationMin', s.durationMin)
    const start = selectedDate()!
    const hours = parseInt(s.startTime.slice(0, 2))
    const minutes = parseInt(s.startTime.slice(3))
    start.setHours(hours, minutes, 0, 0)
    setTraining('start', start)
    setCurrentComponent((c) => c + 1)
  }

  const validateWhenPickingSession = () => {
    let isValid = true
    if (session() === undefined || session() === 'not-selected') {
      setSelectedSession(undefined)
      isValid = false
    }

    return isValid
  }

  const validateWhenSettingSession = () => {
    let isValid = true
    if (day() === undefined || day() === NullDay) {
      setDay(undefined)
      isValid = false
    }

    if (startTime() === undefined || startTime() === NullStartTime) {
      setStartTime(undefined)
      isValid = false
    }

    if (startTime()?.match(/^[0-9]{2}:[0-9]{2}$/) === null) {
      isValid = false
    }
    return isValid
  }

  const manualTrainingSessionForm = () => {
    return (
      <div class="m-4">
        <div class="flex items-center justify-between">
          <label class="mx-4 text-xl font-bold" for="day">
            <Trans key="day" />
          </label>
          <select
            id="day"
            classList={{
              'border-red-500': day() === undefined,
              'border-slate-300': day() !== undefined
            }}
            class="my-2 rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
            onChange={(e) => {
              setDay(e.target.value as Day)
              setSelectedDate(NullDateTime)
            }}
          >
            <option value="" disabled={day() !== NullDay}>
              <Trans key="select.day" />
            </option>
            <For each={Object.keys(Day)}>
              {(d) => (
                <option selected={d === day()} value={d}>
                  <Trans key={d.toLowerCase()} />
                </option>
              )}
            </For>
          </select>
        </div>
        <div class="flex items-center justify-between">
          <label class="mx-4 w-3/4 text-xl font-bold" for="duration">
            <Trans key="duration.in.minutes" />
          </label>
          <input
            id="duration"
            type="number"
            min="1"
            placeholder={'60'}
            classList={{
              'border-red-500 text-red-500': durationMin() === undefined,
              'border-slate-300': durationMin() !== undefined
            }}
            class="my-2 w-1/4 rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
            value={durationMin()}
            onChange={(e) => {
              const val = e.target.value
              const dur = parseInt(val)
              if (Number.isNaN(dur) || dur < 1 || dur > SmallIntMax) {
                setDurationMin(undefined)
                return
              }
              setDurationMin(dur)
            }}
          />
        </div>
        <div class="flex items-center justify-between">
          <label class="mx-4 text-xl font-bold">
            <Trans key="starttime" />
          </label>
          <div
            classList={{
              'border-red-500 text-red-500': startTime() === undefined,
              'border-slate-300': startTime() !== undefined
            }}
            class="my-2 w-auto rounded-md border px-4 py-2 text-xl"
          >
            <div class="flex">
              <select
                name="hours"
                class="appearance-none bg-transparent text-xl outline-none"
                onChange={(e) => {
                  const val = e.target.value
                  const minutes = startTime() ? startTime()!.slice(3) : '--'
                  const newTime = val + ':' + minutes
                  setStartTime(newTime)
                }}
              >
                <option
                  value=""
                  disabled={startTime() !== NullStartTime.slice(0, 2)}
                >
                  --
                </option>
                <For each={Array.from(StartTimeHours)}>
                  {(hour) => (
                    <option
                      selected={
                        startTime() !== undefined &&
                        hour === startTime()!.slice(0, 2)
                      }
                      value={hour}
                    >
                      {hour}
                    </option>
                  )}
                </For>
              </select>
              <span class="mx-2">:</span>
              <select
                name="minutes"
                class="appearance-none bg-transparent text-xl outline-none"
                onChange={(e) => {
                  const val = e.target.value
                  const hours = startTime() ? startTime()!.slice(0, 2) : '--'
                  const newTime = hours + ':' + val
                  setStartTime(newTime)
                }}
              >
                <option
                  value=""
                  disabled={startTime() !== NullStartTime.slice(3)}
                >
                  --
                </option>
                <For each={Array.from(StartTimeMinutes)}>
                  {(minute) => (
                    <option
                      selected={
                        startTime() !== undefined &&
                        minute === startTime()!.slice(3)
                      }
                      value={minute}
                    >
                      {minute}
                    </option>
                  )}
                </For>
              </select>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div class="h-screen">
      <div class="h-3/5">
        <Show when={pickSession()} fallback={manualTrainingSessionForm()}>
          <Switch
            fallback={
              <SessionPicker
                sessions={sessions()}
                selectedSession={session()}
                sessionPage={sessionsPage()}
                isLastPage={isLastPage()}
                onNextPage={() => {
                  setSelectedSession('not-selected')
                  fetchNextSessionPage()
                }}
                onPrevPage={() => {
                  setSelectedSession('not-selected')
                  fetchPrevSessionPage()
                }}
                onSelect={(s) => {
                  setSelectedSession(s)
                  setSelectedDate(NullDateTime)
                }}
              />
            }
          >
            <Match when={sessions()?.length === 0}>
              <div class="m-4 rounded bg-blue-200 p-4 text-lg font-bold">
                <Trans key="no.sessions.created" />
                <button class="my-4 block rounded-lg bg-sky-300 p-2 text-xl font-bold text-black">
                  <A href="/session/create">
                    <Trans key="create.session" />
                  </A>
                </button>
              </div>
            </Match>
          </Switch>
        </Show>
      </div>
      <div class="h-2/5">
        <div class="flex w-5/6 items-center text-center">
          <label class="mx-4 w-1/2 text-xl font-bold" for="date">
            <Trans key="date" />
          </label>
          <select
            id="date"
            disabled={session() === 'not-selected' || session() === undefined}
            classList={{
              'border-red-500 text-red': selectedDate() === undefined,
              'border-slate-300 text-black': selectedDate() !== undefined
            }}
            class="my-2 w-11/12 rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
            onChange={(e) => {
              const date = new Date(dates()[parseInt(e.target.value)])
              setSelectedDate(date)
            }}
          >
            <option value="" disabled={training.start !== NullDateTime}>
              <Trans key="select.date" />
            </option>
            <For each={dates()}>
              {(d, i) => {
                return (
                  <option
                    selected={
                      training.start !== NullDateTime &&
                      training.start.toDateString() ===
                        dates()[i()].toDateString()
                    }
                    value={i()}
                  >
                    {formatDate(d)}
                  </option>
                )
              }}
            </For>
          </select>
        </div>
        <div class="flex w-screen justify-around text-xl">
          <button
            classList={{
              'bg-sky-500 text-white': pickSession(),
              'bg-white text-sky-500': !pickSession(),
              'cursor-not-allowed opacity-50': serverError() !== undefined
            }}
            disabled={serverError() !== undefined}
            class="rounded border border-sky-500 px-4 py-2 text-xl font-bold"
            onClick={() => setPickSession(true)}
          >
            <Trans key="assign.session" />
          </button>
          <button
            classList={{
              'bg-white text-sky-500': pickSession(),
              'bg-sky-500 text-white': !pickSession()
            }}
            class="rounded border border-sky-500 px-4 py-2 font-bold"
            onClick={() => setPickSession(false)}
          >
            <Trans key="set.manually" />
          </button>
        </div>
      </div>

      <button
        class="fixed bottom-0 right-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
        onClick={() => sumbit()}
      >
        <Trans key="next" />
      </button>
    </div>
  )
}
