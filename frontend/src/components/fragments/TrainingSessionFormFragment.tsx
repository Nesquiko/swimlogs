import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { Component, For, Show } from 'solid-js'
import { Day, Session } from '../../generated'
import {
  NullDate,
  NullDay,
  NullStartTime,
  SmallIntMax,
  StartTimeHours,
  StartTimeMinutes
} from '../../lib/consts'
import { formatDate } from '../../lib/datetime'
import { useCreateTraining } from '../context/CreateTrainingContextProvider'

export const TrainingSessionForm: Component = () => {
  const [
    { training, setTraining },
    { invalidTraining, setInvalidTraining },
    sessions,
    state,
    [, setCurrentComponent]
  ] = useCreateTraining()
  const [fromSession, setFromSession] = state.fromSession
  const [selectedSession, setSelectedSession] = state.selectedSession
  const [day, setDay] = state.day
  const [dates] = state.dates
  const [selectedDate, setSelectedDate] = state.selectedDate
  const [t] = useTransContext()

  const sumbit = () => {
    let isValid = true
    if (fromSession() && selectedSession() === '') {
      setSelectedSession(undefined)
      isValid = false
    }
    if (!fromSession() && day() === NullDay) {
      setDay(undefined)
      isValid = false
    }
    if (!fromSession() && invalidTraining.durationMin !== undefined) {
      isValid = false
    }
    if (!fromSession() && training.startTime.includes('--')) {
      setInvalidTraining('startTime', 'Start time must be selected')
      isValid = false
    }
    if (training.date === NullDate) {
      setSelectedDate(undefined)
      isValid = false
    }
    if (!isValid) {
      return
    }
    setCurrentComponent((c) => c + 1)
  }

  const manualTrainingSessionForm = () => {
    return (
      <div class="m-4">
        <select
          classList={{
            'border-red-500': day() === undefined,
            'border-slate-300': day() !== undefined
          }}
          class="my-2 w-full rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            setDay(e.target.value as Day)
            setTraining('date', NullDate)
            setSelectedDate(NullDate)
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
        <div class="flex items-center">
          <label class="mx-4 w-3/4 text-xl font-bold" for="duration">
            <Trans key="duration.in.minutes" />
          </label>
          <input
            id="duration"
            type="number"
            min="1"
            placeholder={t('minutes', 'minutes')}
            classList={{
              'border-red-500 text-red-500':
                invalidTraining.durationMin !== undefined,
              'border-slate-300': invalidTraining.durationMin === undefined
            }}
            class="w-1/4 rounded-md border px-4 py-2 text-xl focus:border-blue-500 focus:outline-none focus:ring"
            value={training.durationMin}
            onChange={(e) => {
              const val = e.target.value
              setInvalidTraining('durationMin', undefined)
              const dur = parseInt(val)
              if (Number.isNaN(dur) || dur < 1 || dur > SmallIntMax) {
                setInvalidTraining(
                  'durationMin',
                  'Duration must be between 1 and 32767'
                )
                return
              }

              setTraining('durationMin', dur)
            }}
          />
        </div>
        <div class="flex items-center justify-between">
          <label class="mx-4 text-xl font-bold">
            <Trans key="starttime" />
          </label>
          <div
            classList={{
              'border-red-500 text-red-500':
                invalidTraining.startTime !== undefined,
              'border-slate-300': invalidTraining.startTime === undefined
            }}
            class="my-2 w-auto rounded-md border px-4 py-2 text-xl"
          >
            <div class="flex">
              <select
                name="hours"
                class="appearance-none bg-transparent text-xl outline-none"
                onChange={(e) => {
                  const val = e.target.value
                  const newTime = val + ':' + training.startTime.slice(3)
                  setInvalidTraining('startTime', undefined)
                  setTraining('startTime', newTime)
                }}
              >
                <option
                  value=""
                  disabled={training.startTime !== NullStartTime.slice(0, 2)}
                >
                  --
                </option>
                <For each={Array.from(StartTimeHours)}>
                  {(hour) => (
                    <option
                      selected={hour === training.startTime.slice(0, 2)}
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
                  const newTime = training.startTime.slice(0, 3) + val
                  setInvalidTraining('startTime', undefined)
                  setTraining('startTime', newTime)
                }}
              >
                <option
                  value=""
                  disabled={training.startTime !== NullStartTime.slice(3)}
                >
                  --
                </option>
                <For each={Array.from(StartTimeMinutes)}>
                  {(minute) => (
                    <option
                      selected={minute === training.startTime.slice(3)}
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

  const setFromSessionForm = () => {
    return (
      <div class="m-4">
        <select
          classList={{
            'border-red-500': selectedSession() === undefined,
            'border-slate-300': selectedSession() !== undefined
          }}
          class="my-2 w-full rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            const sessionStr = e.target.value
            const session = JSON.parse(sessionStr) as Session

            const localizedDay = t(session.day.toLowerCase(), session.day)
            const sesName = `${localizedDay} ${session.startTime} (${session.durationMin} min)`
            setSelectedSession(sesName)

            setDay(session.day)
            setTraining('startTime', session.startTime)
            setTraining('durationMin', session.durationMin)
            setTraining('date', NullDate)
            setSelectedDate(NullDate)
          }}
        >
          <option value="" disabled={selectedSession() !== ''}>
            <Trans key="select.session" />
          </option>
          <For each={sessions()?.sessions ?? []}>
            {(session) => {
              const localizedDay = t(session.day.toLowerCase(), session.day)
              const sesName = `${localizedDay} ${session.startTime} (${session.durationMin} min)`
              const sessionStr = JSON.stringify(session)
              return (
                <option
                  selected={selectedSession() === sesName}
                  value={sessionStr}
                >
                  {sesName}
                </option>
              )
            }}
          </For>
        </select>
      </div>
    )
  }

  return (
    <div>
      <Show when={fromSession()} fallback={manualTrainingSessionForm()}>
        {setFromSessionForm()}
      </Show>
      <div class="fixed bottom-1/2 flex flex-col items-center justify-center">
        <div class="flex w-5/6 items-center text-center">
          <label class="mx-4 w-1/2 text-xl font-bold" for="date">
            <Trans key="date" />
          </label>
          <select
            id="date"
            classList={{
              'border-red-500 text-red': selectedDate() === undefined,
              'border-slate-300 text-black': selectedDate() !== undefined
            }}
            class="my-2 w-11/12 rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
            onChange={(e) => {
              const date = dates()[parseInt(e.target.value)]
              setSelectedDate(date)
              setTraining('date', date)
            }}
          >
            <option value="" disabled={training.date !== NullDate}>
              <Trans key="select.date" />
            </option>
            <For each={dates()}>
              {(d, i) => {
                return (
                  <option
                    selected={
                      training.date === dates()[i()] &&
                      training.date !== NullDate
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
        <div class="my-4 flex w-screen justify-around text-xl">
          <button
            classList={{
              'bg-white text-sky-500': fromSession(),
              'bg-sky-500 text-white': !fromSession()
            }}
            class="rounded border border-sky-500 px-4 py-2 font-bold"
            onClick={() => setFromSession(false)}
          >
            <Trans key="set.manually" />
          </button>
          <button
            classList={{
              'bg-sky-500 text-white': fromSession(),
              'bg-white text-sky-500': !fromSession(),
              'cursor-not-allowed opacity-50': sessions()?.sessions.length === 0
            }}
            class="rounded border border-sky-500 px-4 py-2 text-xl font-bold"
            disabled={sessions()?.sessions.length === 0}
            onClick={() => setFromSession(true)}
          >
            <Trans key="assign.session" />
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
