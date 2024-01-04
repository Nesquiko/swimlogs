import {
  Accessor,
  createContext,
  ParentComponent,
  Resource,
  useContext,
} from 'solid-js'
import { Session } from '../generated'

const makeSessionsContext = (
  sessions: Resource<Session[]>,
  fetchNextSessionPage: () => void,
  fetchPrevSessionPage: () => void,
  page: Accessor<number>,
  isLastPage: Accessor<boolean>,
  serverError: Accessor<string | undefined>
) => {
  return [
    sessions,
    fetchNextSessionPage,
    fetchPrevSessionPage,
    page,
    isLastPage,
    serverError,
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
  serverError: Accessor<string | undefined>
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
        props.isLastPage,
        props.serverError
      )}
    >
      {props.children}
    </SessionsContext.Provider>
  )
}
