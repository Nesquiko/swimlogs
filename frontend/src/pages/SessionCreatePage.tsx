import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { useNavigate } from '@solidjs/router'
import { Component, createSignal, For } from 'solid-js'
import { createStore } from 'solid-js/store'
import { openToast, ToastType } from '../components/Toast'
import { CreateSessionRequest, Day, ResponseError } from '../generated'
import { SmallIntMax, StartTimeHours, StartTimeMinutes } from '../lib/consts'
import { sessionApi } from '../state/session'

const SessionCreatePage: Component = () => {
  const [day, setDay] = createSignal<Day>(Day.Monday)
  const [duration, setDuration] = createSignal<number>(60)
  const [startTime, setStartTime] = createStore<{
    hours: number
    minutes: number
  }>({ hours: 6, minutes: 0 })

  const [t] = useTransContext()
  const navigate = useNavigate()
  async function createSession() {
    const session: CreateSessionRequest = {
      startTime: `${startTime.hours
        .toString()
        .padStart(2, '0')}:${startTime.minutes.toString().padStart(2, '0')}`,
      durationMin: duration(),
      day: day()
    }
    await sessionApi
      .createSession({
        createSessionRequest: session
      })
      .then(() => {
        openToast(t('session.created', 'Session created'), ToastType.SUCCESS)
        navigate('/', { replace: true })
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        openToast(
          t('session.creation.error', 'Error creating session'),
          ToastType.ERROR
        )
        navigate('/', { replace: true })
      })
  }

  return (
    <div>
      <h1 class="m-4 text-2xl font-bold">
        <Trans key="create.session" />
      </h1>
      <form class="mx-4">
        <fieldset>
          <label class="block text-xl" for="day">
            <Trans key="day" />
          </label>
          <select
            id="day"
            class="text-cente w-32 rounded-lg border border-solid border-slate-300 bg-white p-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
            onChange={(e) => setDay(e.target.value as Day)}
          >
            <For each={Object.keys(Day)}>
              {(d) => (
                <option selected={d === day()} value={d}>
                  <Trans key={d.toLowerCase()} />
                </option>
              )}
            </For>
          </select>
          <label class="mt-4 block text-xl" for="duration">
            <Trans key="duration.in.minutes" />
          </label>
          <input
            id="duration"
            type="number"
            min="1"
            max={SmallIntMax}
            class="w-32 rounded-lg border border-solid border-slate-300 bg-white p-2 text-center text-xl focus:border-sky-500 focus:outline-none focus:ring"
            value={duration()}
            onChange={(e) => {
              const val = e.target.value
              const dur = parseInt(val)
              if (Number.isNaN(dur) || dur < 1 || dur > SmallIntMax) {
                setDuration(60)
                return
              }
              setDuration(dur)
            }}
          />
          <label class="mt-4 block text-xl" for="hours">
            <Trans key="starttime" />
          </label>
          <div class="w-32 rounded-lg border border-slate-300 px-4 py-2 text-center text-xl">
            <select
              id="hours"
              class="appearance-none bg-transparent outline-none"
              onChange={(e) => {
                const val = e.target.value
                const hours = parseInt(val)
                setStartTime('hours', hours)
              }}
            >
              <For each={Array.from(StartTimeHours)}>
                {(hour) => (
                  <option
                    selected={parseInt(hour) === startTime.hours}
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
              class="appearance-none bg-transparent outline-none"
              onChange={(e) => {
                const val = e.target.value
                const minutes = parseInt(val)
                setStartTime('minutes', minutes)
              }}
            >
              <For each={Array.from(StartTimeMinutes)}>
                {(minute) => (
                  <option
                    selected={parseInt(minute) === startTime.minutes}
                    value={minute}
                  >
                    {minute}
                  </option>
                )}
              </For>
            </select>
          </div>
        </fieldset>
      </form>
      <button
        class="m-4 rounded-lg bg-sky-500 p-4 text-xl font-bold text-white"
        onClick={createSession}
      >
        <Trans key="create.session" />
      </button>
    </div>
  )
}

export default SessionCreatePage
