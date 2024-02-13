import 'flowbite'
import { Component, createSignal } from 'solid-js'
import { Drawer } from './components/Drawer'
import Header from './components/Header'
import DismissibleToast, {
  ToastMode,
} from './components/common/DismissibleToast'
import { Route, Router } from '@solidjs/router'
import { TransProvider } from '@mbarzda/solid-i18next'
import i18next from 'i18next'
import I18NextHttpBackend from 'i18next-http-backend'
import I18nextBrowserLanguageDetector from 'i18next-browser-languagedetector'
import Home from './pages/Home'
import TrainingPage from './pages/TrainingPage'
import TrainingHistoryPage from './pages/TrainingsHistoryPage'
import NewTrainingPage from './pages/NewTrainingPage'

const [openToast, setOpenToast] = createSignal(false)
const [toastMessage, setToastMessage] = createSignal('')
const [toastMode, setToastMode] = createSignal(ToastMode.SUCCESS)

const showToast = (
  message: string,
  mode: ToastMode = ToastMode.SUCCESS
): void => {
  setToastMessage(message)
  setToastMode(mode)
  setOpenToast(true)
}

const App: Component = () => {
  i18next.use(I18NextHttpBackend)
  i18next.use(I18nextBrowserLanguageDetector)
  const backend = { loadPath: '/locales/{{lng}}/{{ns}}.json' }

  return (
    <Router
      root={(props) => (
        <TransProvider
          options={{
            backend,
            fallbackLng: 'en',
          }}
        >
          <Header />
          <Drawer />
          <DismissibleToast
            open={openToast()}
            onDismiss={() => setOpenToast(false)}
            mode={toastMode()}
            message={toastMessage()}
          />
          <div class="py-2">{props.children}</div>
        </TransProvider>
      )}
    >
      <Route path="/" component={Home} />
      <Route path="/training/new" component={NewTrainingPage} />
      <Route path="/training/:id" component={TrainingPage} />
      <Route path="/trainings" component={TrainingHistoryPage} />
    </Router>
  )
}

export default App
export { showToast }
