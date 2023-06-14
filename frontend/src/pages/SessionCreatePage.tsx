import { useNavigate } from '@solidjs/router'
import { Component, createSignal, For } from 'solid-js'
import { createStore } from 'solid-js/store'
import { openToast, ToastType } from '../components/Toast'
import { CreateSessionRequest, Day, ResponseError } from '../generated'
import { SmallIntMax, StartTimeHours, StartTimeMinutes } from '../lib/consts'
import { sessionApi } from '../state/session'

const SessionCreatePage: Component = () => {
  const [session, setSession] = createStore<CreateSessionRequest>({
    day: '' as Day,
    durationMin: 60,
    startTime: ''
  })
  const [startTime, setStartTime] = createStore({ hours: '06', minutes: '00' })
  const [isDaySelected, setIsDaySelected] = createSignal(false)
  const [dayError, setDayError] = createSignal<boolean>(false)
  const [durationError, setDurationError] = createSignal<boolean>(false)

  const navigate = useNavigate()
  async function createSession() {
    const starT = startTime.hours + ':' + startTime.minutes
    setSession('startTime', starT)

    if (!checkInputErrors()) return

    await sessionApi
      .createSession({
        createSessionRequest: session
      })
      .then(() => {
        openToast('Session created', ToastType.SUCCESS)
        navigate('/', { replace: true })
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        openToast('Error creating session', ToastType.ERROR)
        navigate('/', { replace: true })
      })
  }

  const checkInputErrors = () => {
    let isValid = true

    if (!isDaySelected()) {
      setDayError(true)
      isValid = false
    } else {
      setDayError(false)
    }

    if (durationError()) {
      isValid = false
    }

    return isValid
  }

  const manualTrainingSessionForm = () => {
    return (
      <div class="m-4">
        <select
          classList={{
            'border-red-500': dayError(),
            'border-slate-300': !dayError()
          }}
          class="my-2 w-full rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
          onChange={(e) => {
            setIsDaySelected(true)
            setDayError(false)
            const val = e.target.value
            setSession('day', val as Day)
          }}
        >
          <option value="" disabled={isDaySelected()}>
            Select day
          </option>
          <For each={Object.keys(Day)}>
            {(day) => <option value={day}>{day}</option>}
          </For>
        </select>

        <div class="flex items-center">
          <label class="mx-4 w-3/4 text-xl font-bold" for="duration">
            Duration in minutes
          </label>
          <input
            id="duration"
            type="number"
            min="1"
            placeholder="minutes"
            classList={{
              'border-red-500 text-red-500': durationError(),
              'border-slate-300': !durationError()
            }}
            class="w-1/4 rounded-md border px-4 py-2 text-xl focus:border-blue-500 focus:outline-none focus:ring"
            value={session.durationMin}
            onChange={(e) => {
              const val = (e.target as HTMLInputElement).value
              setDurationError(false)
              const dur = parseInt(val)
              if (Number.isNaN(dur) || dur < 1 || dur > SmallIntMax) {
                setDurationError(true)
                return
              }
              setSession('durationMin', dur)
            }}
          />
        </div>

        <div class="flex items-center justify-between">
          <label class="mx-4 text-xl font-bold">Start time</label>
          <div class="my-2 flex w-auto rounded-md border border-slate-300 px-4 py-2 text-xl">
            <select
              name="hours"
              class="appearance-none bg-transparent text-xl outline-none"
              onChange={(e) => {
                const val = e.target.value
                setStartTime('hours', val)
              }}
            >
              <For each={Array.from(StartTimeHours)}>
                {(hour) => <option value={parseInt(hour)}>{hour}</option>}
              </For>
            </select>
            <span class="mx-2">:</span>
            <select
              name="minutes"
              class="appearance-none bg-transparent text-xl outline-none"
              onChange={(e) => {
                const val = e.target.value
                setStartTime('minutes', val)
              }}
            >
              <For each={Array.from(StartTimeMinutes)}>
                {(minute) => <option value={parseInt(minute)}>{minute}</option>}
              </For>
            </select>
          </div>
        </div>
      </div>
    )
  }
  return (
    <div>
      <h1 class="m-4 text-2xl font-bold">Create Session</h1>

      {manualTrainingSessionForm()}

      <div class="text-center">
        <button
          class="mx-auto rounded bg-green-600 p-4 text-xl font-bold text-white"
          onClick={createSession}
        >
          Create Session
        </button>
      </div>
    </div>
  )
}

export default SessionCreatePage
