import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { A } from '@solidjs/router'
import {
  Component,
  createEffect,
  createSignal,
  For,
  Match,
  Switch
} from 'solid-js'
import SessionPicker from '../components/SessionPicker'
import { useSessionsContext } from '../components/SessionsContextProvider'
import { useShownComponent } from '../components/ShownComponentContextProvider'
import { openToast, ToastType } from '../components/Toast'
import { useStateContext } from '../components/TrainingStateContext'
import { Day, Session } from '../generated'
import { SmallIntMax, StartTimeHours, StartTimeMinutes } from '../lib/consts'
import plusSvgBlack from '../assets/plus-black.svg'
import editSvgBlack from '../assets/edit-black.svg'
import selectSvgBlack from '../assets/select-black.svg'
import DatePickerModal from '../components/DatePickerModal'

export const TrainingSessionForm: Component = () => {
  const [t] = useTransContext()
  const [{ setTraining }, state] = useStateContext()
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

  const [pickSession, setPickSession] = state.pickSession
  const [day, setDay] = state.day
  const [durationMin, setDurationMin] = state.durationMin
  const [startTime, setStartTime] = state.startTime
  const [session, setSelectedSession] = state.session

  const [datePickerOpener, setDatepickerOpener] = createSignal({ s: session() })

  const sumbit = (start: Date, s: Session) => {
    setTraining('durationMin', s.durationMin)
    setTraining('start', start)
    setCurrentComponent((c) => c + 1)
  }

  return (
    <div class="h-[32rem]">
      <DatePickerModal opener={datePickerOpener()} onSubmit={sumbit} />
      <p class="m-4 text-2xl font-bold">
        <Trans key="select.session" />
        <Switch>
          <Match when={serverError() !== undefined}>
            <></>
          </Match>
          <Match when={!pickSession()}>
            <img
              class="float-right inline-block cursor-pointer"
              src={selectSvgBlack}
              width={32}
              height={32}
              onClick={() => {
                setPickSession(true)
              }}
            />
          </Match>
          <Match when={pickSession()}>
            <img
              class="float-right inline-block cursor-pointer"
              src={editSvgBlack}
              width={32}
              height={32}
              onClick={() => {
                setPickSession(false)
              }}
            />
          </Match>
        </Switch>
      </p>
      <Switch>
        <Match
          when={
            sessions() !== undefined && sessions()!.length > 0 && pickSession()
          }
        >
          <SessionPicker
            sessions={sessions()}
            selectedSession={session()}
            sessionPage={sessionsPage()}
            isLastPage={isLastPage()}
            onNextPage={() => {
              setSelectedSession((s) => {
                s.id = ''
                return s
              })
              fetchNextSessionPage()
            }}
            onPrevPage={() => {
              setSelectedSession((s) => {
                s.id = ''
                return s
              })
              fetchPrevSessionPage()
            }}
            onSelect={(s) => setSelectedSession(s)}
          />
        </Match>
        <Match when={sessions()?.length === 0 && pickSession()}>
          <div class="text-center">
            <div class="m-4 rounded-lg bg-sky-200 p-4 text-start text-lg font-bold">
              <Trans key="no.sessions.created" />
            </div>
            <button class="text-md m-4 rounded-lg bg-yellow-400 p-2 font-bold text-black">
              <img
                class="mr-2 inline"
                src={plusSvgBlack}
                width={32}
                height={32}
              />
              <A href="/session/create">
                <Trans key="create.session" />
              </A>
            </button>
            <button
              class="text-md m-4 rounded-lg bg-yellow-400 p-2 font-bold text-black"
              onClick={() => setPickSession(false)}
            >
              <img
                class="mr-2 inline"
                src={editSvgBlack}
                width={32}
                height={32}
              />
              <Trans key="set.session" />
            </button>
          </div>
        </Match>
        <Match when={!pickSession()}>
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
                value={durationMin()}
                onChange={(e) => {
                  const val = e.target.value
                  const dur = parseInt(val)
                  if (Number.isNaN(dur) || dur < 1 || dur > SmallIntMax) {
                    setDurationMin(60)
                    return
                  }
                  setDurationMin(dur)
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
                    setStartTime({
                      hours: hours,
                      minutes: startTime().minutes
                    })
                  }}
                >
                  <For each={Array.from(StartTimeHours)}>
                    {(hour) => (
                      <option
                        selected={parseInt(hour) === startTime().hours}
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
                    setStartTime({
                      hours: startTime().hours,
                      minutes: minutes
                    })
                  }}
                >
                  <For each={Array.from(StartTimeMinutes)}>
                    {(minute) => (
                      <option
                        selected={parseInt(minute) === startTime().minutes}
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
        </Match>
      </Switch>
      <button
        class="absolute bottom-0 right-4 my-4 w-20 rounded-lg bg-sky-500 py-2 text-xl font-bold text-white"
        onClick={() => {
          if (pickSession() && session().id === '') {
            openToast(t('please.select.session'))
            return
          }
          setDatepickerOpener({ s: session() })
        }}
      >
        <Trans key="next" />
      </button>
    </div>
  )
}
