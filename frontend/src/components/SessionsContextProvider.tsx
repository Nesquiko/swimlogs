import {
  Accessor,
  createContext,
  ParentComponent,
  Resource,
  useContext
} from 'solid-js'
import { Session } from '../generated'

const makeSessionsContext = (
  sessions: Resource<Session[]>,
  fetchNextSessionPage: () => void,
  fetchPrevSessionPage: () => void,
  page: Accessor<number>,
  isLastPage: Accessor<boolean>
) => {
  return [
    sessions,
    fetchNextSessionPage,
    fetchPrevSessionPage,
    page,
    isLastPage
  ] as const
}

type SessionsContextType = ReturnType<typeof makeSessionsContext>

const SessionsContext = createContext<SessionsContextType>()

export const useSessionsContext = () => useContext(SessionsContext)!

interface SessionsContextProps {
  sessions: Resource<Session[]>
  fetchNextSessionPage: () => void
  fetchPrevSessionPage: () => void
  page: Accessor<number>
  isLastPage: Accessor<boolean>
}

export const SessionsContextProvider: ParentComponent<SessionsContextProps> = (
  props
) => {
  return (
    <SessionsContext.Provider
      value={makeSessionsContext(
        props.sessions,
        props.fetchNextSessionPage,
        props.fetchPrevSessionPage,
        props.page,
        props.isLastPage
      )}
    >
      {props.children}
    </SessionsContext.Provider>
  )
}
