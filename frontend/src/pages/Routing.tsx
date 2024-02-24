import { Route, useLocation } from '@solidjs/router'
import { Component, createContext, ParentComponent, useContext } from 'solid-js'
import Home from './Home'
import NewTrainingPage from './NewTrainingPage'
import TrainingHistoryPage from './TrainingsHistoryPage'
import TrainingDisplayPage from './TrainingDisplayPage'
import { loadTrainingById } from '../state/trainings'
import TrainingEditPage from './TrainingEditPage'

const makeOnBackContext = (
  onBack: () => void,
  setOnBack: (onBack?: () => void) => void
) => {
  return [onBack, setOnBack] as const
}
type OnBackContextType = ReturnType<typeof makeOnBackContext>
const OnBackContext = createContext<OnBackContextType>()

export const useOnBackcontext = () => useContext(OnBackContext)!

export const OnBackContextProvider: ParentComponent = (props) => {
  const location = useLocation()
  let onBackOverride: (() => void) | undefined = undefined

  const back = () => {
    if (onBackOverride) {
      onBackOverride()
      return
    }

    if (location.pathname === '/') {
      return
    }
    history.back()
  }

  return (
    <OnBackContext.Provider
      value={makeOnBackContext(back, (onBack) => {
        onBackOverride = onBack
      })}
    >
      {props.children}
    </OnBackContext.Provider>
  )
}

const Routes: Component = () => {
  return (
    <>
      <Route path="/" component={Home} />
      <Route path="/training/new" component={NewTrainingPage} />
      {/* <Route path="/training/:id" component={EditTrainingPage}> */}
      <Route path="/training/:id">
        <Route
          path="/display"
          component={TrainingDisplayPage}
          load={(args) => {
            return { trainingPromise: loadTrainingById(args) }
          }}
        />
        <Route
          path="/edit"
          component={TrainingEditPage}
          load={(args) => {
            return { trainingPromise: loadTrainingById(args) }
          }}
        />
        <Route path="/edit/session" component={() => <>Edit Session</>} />
      </Route>

      <Route path="/trainings" component={TrainingHistoryPage} />
    </>
  )
}

export default Routes
