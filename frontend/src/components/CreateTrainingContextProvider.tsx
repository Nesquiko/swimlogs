import {
  createContext,
  createEffect,
  createResource,
  createSignal,
  ParentComponent,
  Resource,
  Signal,
  useContext
} from 'solid-js'
import { createStore } from 'solid-js/store'
import {
  Day,
  GetSessionsResponse,
  InvalidTraining,
  InvalidTrainingSet,
  NewTraining,
  ResponseError,
  Session
} from '../generated'
import { NullDateTime, NullDay, NullStartTime, PAGE_SIZE } from '../lib/consts'
import { sessionApi } from '../state/session'

const makeTrainingContext = (
  newTraining: NewTraining,
  sessions: Resource<GetSessionsResponse>,
  state: State,
  currentComponentSignal: Signal<number>,
  sumbitTraining: () => void
) => {
  const [training, setTraining] = createStore<NewTraining>(newTraining)
  const [invalidTraining, setInvalidTraining] = createStore<InvalidTraining>({
    invalidSets: training.sets.map(() => {
      return {} as InvalidTrainingSet
    })
  })

  return [
    { training, setTraining },
    { invalidTraining, setInvalidTraining },
    sessions,
    state,
    [currentComponentSignal[0], currentComponentSignal[1]],
    sumbitTraining
  ] as const
}
type CreateTrainingContextType = ReturnType<typeof makeTrainingContext>

const CreateTrainingContext = createContext<CreateTrainingContextType>()

export const useCreateTraining = () => useContext(CreateTrainingContext)!

interface CreateTrainingContextProps {
  newTraining: NewTraining
  currentComponentSignal: Signal<number>
  sumbitTraining: () => void
}

// TODO split this big context into smaller ones
export const CreateTrainingContextProvider: ParentComponent<
  CreateTrainingContextProps
> = (props) => {
  const state = initialState()

  const [sessionsPage] = state.sessionsPage
  const [, setTotalSessions] = state.totalSessions
  const [session, setSession] = state.session
  const [day] = state.day
  const [durationMin] = state.durationMin
  const [startTime] = state.startTime
  const [sessions] = createResource(sessionsPage, getSessions)
  const [, setFromSession] = state.pickSession
  const [, setDates] = state.dates

  async function getSessions(page: number): Promise<GetSessionsResponse> {
    return sessionApi
      .getSessions({ page, pageSize: PAGE_SIZE })
      .then((res) => {
        setTotalSessions(res.pagination.total)
        return res
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        return Promise.resolve({
          sessions: [],
          pagination: { page: 0, pageSize: PAGE_SIZE, total: 0 }
        })
      })
  }

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
    setFromSession(sessions()?.sessions?.length !== 0)
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
    <CreateTrainingContext.Provider
      value={makeTrainingContext(
        props.newTraining,
        sessions,
        state,
        props.currentComponentSignal,
        props.sumbitTraining
      )}
    >
      {props.children}
    </CreateTrainingContext.Provider>
  )
}

type State = {
  pickSession: Signal<boolean>

  // TODO somehow organize this? maybe store?
  day: Signal<Day | undefined>
  durationMin: Signal<number | undefined>
  startTime: Signal<string | undefined>

  dates: Signal<Array<Date>>
  selectedDate: Signal<Date | undefined>

  sessionsPage: Signal<number>
  totalSessions: Signal<number>

  session: Signal<Session | 'not-selected' | undefined>
}

function initialState(): State {
  const pickSession = createSignal(false)
  const day = createSignal<Day | undefined>(NullDay)
  const durationMin = createSignal<number | undefined>(60)
  const startTime = createSignal<string | undefined>(NullStartTime)
  const dates = createSignal<Array<Date>>([])
  const selectedDate = createSignal<Date | undefined>(NullDateTime)
  const sessionsPage = createSignal(0)
  const totalSessions = createSignal(0)
  const session = createSignal<Session | 'not-selected' | undefined>(
    'not-selected'
  )

  return {
    pickSession,
    day,
    durationMin,
    startTime,
    dates,
    selectedDate,
    sessionsPage,
    totalSessions,
    session
  }
}
