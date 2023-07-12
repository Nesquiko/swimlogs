import {
  createContext,
  createEffect,
  createSignal,
  ParentComponent,
  Signal,
  useContext
} from 'solid-js'
import { createStore } from 'solid-js/store'
import { Day, NewTraining, Session } from '../generated'
import { NullDateTime, NullDay, NullStartTime } from '../lib/consts'

type State = {
  day: Signal<Day | undefined>
  durationMin: Signal<number | undefined>
  startTime: Signal<string | undefined>

  dates: Signal<Array<Date>>
  selectedDate: Signal<Date | undefined>

  session: Signal<Session | 'not-selected' | undefined>
}

const makeTrainingStateContext = (newTraining: NewTraining, state: State) => {
  const [training, setTraining] = createStore<NewTraining>(newTraining)
  return [{ training, setTraining }, state] as const
}
type TrainingStateContextType = ReturnType<typeof makeTrainingStateContext>

const TrainingStateContext = createContext<TrainingStateContextType>()

export const useStateContext = () => useContext(TrainingStateContext)!

interface TrainingStateContextProps {
  newTraining: NewTraining
}

export const TrainingStateContextProvider: ParentComponent<
  TrainingStateContextProps
> = (props) => {
  const state = initialState()

  const [session, setSession] = state.session
  const [day] = state.day
  const [durationMin] = state.durationMin
  const [startTime] = state.startTime
  const [, setDates] = state.dates

  createEffect(() => {
    if (day() === undefined || day() === NullDay) {
      return
    }

    if (durationMin() === undefined || durationMin() === 0) {
      return
    }

    if (startTime() === undefined || startTime() === NullStartTime) {
      return
    }

    if (startTime()!.match(/^[0-9]{2}:[0-9]{2}$/) === null) {
      return
    }

    const newSession = {
      id: '',
      day: day()!,
      durationMin: durationMin()!,
      startTime: startTime()!
    }

    setSession(newSession)
  })

  createEffect(() => {
    if (session() === undefined || session() === 'not-selected') {
      return
    }

    const dayNames = Array.from(Object.keys(Day))
    const inputDayIndex = dayNames.indexOf((session() as Session).day)
    const today = new Date()
    const todayIndex = today.getDay() - 1
    const dayDifference = (inputDayIndex - todayIndex + 7) % 7
    const futureDates: Date[] = []
    // Start from today and find the next four dates on the input day
    for (let i = 0; i < 4; i++) {
      const futureDate = new Date(
        today.getTime() + (dayDifference + 7 * i) * 24 * 60 * 60 * 1000
      )
      futureDate.setHours(0, 0, 0, 0)
      futureDates.push(futureDate)
    }

    setDates(futureDates)
  })

  return (
    <TrainingStateContext.Provider
      value={makeTrainingStateContext(props.newTraining, state)}
    >
      {props.children}
    </TrainingStateContext.Provider>
  )
}

function initialState(): State {
  const day = createSignal<Day | undefined>(NullDay)
  const durationMin = createSignal<number | undefined>(60)
  const startTime = createSignal<string | undefined>(NullStartTime)
  const dates = createSignal<Array<Date>>([])
  const selectedDate = createSignal<Date | undefined>(NullDateTime)
  const session = createSignal<Session | 'not-selected' | undefined>(
    'not-selected'
  )

  return {
    day,
    durationMin,
    startTime,
    dates,
    selectedDate,
    session
  }
}
