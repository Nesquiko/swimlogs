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
} from '../../generated'
import { NullDate, NullDay } from '../../lib/consts'
import { sessionApi } from '../../state/session'

const makeTrainingContext = (
  newTraining: NewTraining,
  sessions: Resource<GetSessionsResponse>,
  state: State,
  currentComponentSignal: Signal<number>,
  sumbitTraining: () => void
) => {
  const [training, setTraining] = createStore<NewTraining>(newTraining)
  const [invalidTraining, setInvalidTraining] = createStore<InvalidTraining>({
    blocks: training.blocks.map((b) => {
      return {
        num: b.num,
        sets: b.sets.map((s) => {
          return { num: s.num, startingRule: {} } as InvalidTrainingSet
        })
      }
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

export const PAGE_SIZE = 7

interface CreateTrainingContextProps {
  newTraining: NewTraining
  currentComponentSignal: Signal<number>
  sumbitTraining: () => void
}

export const CreateTrainingContextProvider: ParentComponent<
  CreateTrainingContextProps
> = (props) => {
  const state = initialState()

  const [sessionsPage] = state.sessionsPage
  const [, setTotalSessions] = state.totalSessions
  const [sessions] = createResource(sessionsPage, getSessions)

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
    const [, setFromSession] = state.fromSession
    setFromSession(sessions()?.sessions?.length !== 0)
  })

  createEffect(() => {
    const [day] = state.day
    if (day() === undefined || day() === NullDay) {
      return
    }
    const [, setDates] = state.dates
    const dayNames = Array.from(Object.keys(Day))
    const inputDayIndex = dayNames.indexOf(day()!)
    const today = new Date()
    // Find the index of today's day
    const todayIndex = today.getDay() - 1
    // Calculate the number of days between today and the input day
    const dayDifference = (inputDayIndex - todayIndex + 7) % 7
    const futureDates: Date[] = []
    // Start from today and find the next four dates on the input day
    for (let i = 0; i < 4; i++) {
      const futureDate = new Date(
        today.getTime() + (dayDifference + 7 * i) * 24 * 60 * 60 * 1000
      )
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
  fromSession: Signal<boolean>
  day: Signal<Day | undefined>
  dates: Signal<Array<Date>>
  selectedDate: Signal<Date | undefined>
  currentBlock: Signal<number>
  sessionsPage: Signal<number>
  totalSessions: Signal<number>
  selectedSession: Signal<Session | 'not-selected' | undefined>
}

function initialState(): State {
  const fromSession = createSignal(false)
  const day = createSignal<Day | undefined>(NullDay)
  const dates = createSignal<Array<Date>>([])
  const selectedDate = createSignal<Date | undefined>(NullDate)
  const currentBlock = createSignal(0)
  const sessionsPage = createSignal(0)
  const totalSessions = createSignal(0)
  const selectedSession = createSignal<Session | 'not-selected' | undefined>(
    'not-selected'
  )

  return {
    fromSession,
    day,
    dates,
    selectedDate,
    currentBlock,
    sessionsPage,
    totalSessions,
    selectedSession
  }
}
