import {
  createContext,
  createEffect,
  createSignal,
  ParentComponent,
  Signal,
  useContext,
} from 'solid-js'
import { createStore } from 'solid-js/store'
import { Day, NewTraining, Session } from '../generated'
import { NullDateTime } from '../lib/consts'

type State = {
  day: Signal<Day>
  durationMin: Signal<number>
  startTime: Signal<{ hours: number; minutes: number }>

  selectedDate: Signal<Date | undefined>

  pickSession: Signal<boolean>
  session: Signal<Session>
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

  const [, setSession] = state.session
  const [day] = state.day
  const [durationMin] = state.durationMin
  const [startTime] = state.startTime

  createEffect(() => {
    if (durationMin() <= 0) {
      return
    }

    const startTimeStr = `${startTime()
      .hours.toString()
      .padStart(2, '0')}:${startTime().minutes.toString().padStart(2, '0')}`

    const newSession = {
      id: '',
      day: day()!,
      durationMin: durationMin()!,
      startTime: startTimeStr,
    }

    setSession(newSession)
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
  const day = createSignal<Day>(Day.Monday)
  const durationMin = createSignal<number>(60)
  const startTime = createSignal<{ hours: number; minutes: number }>({
    hours: 6,
    minutes: 0,
  })
  const selectedDate = createSignal<Date | undefined>(NullDateTime)
  const pickSession = createSignal<boolean>(true)
  const session = createSignal<Session>({
    id: '',
    day: Day.Monday,
    durationMin: 60,
    startTime: '06:00',
  })

  return {
    day,
    durationMin,
    startTime,
    selectedDate,
    pickSession,
    session,
  }
}
