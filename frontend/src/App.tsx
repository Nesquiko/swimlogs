import 'flowbite'
import { Routes, Route } from '@solidjs/router'
import { Component, createSignal, lazy } from 'solid-js'
import { Drawer } from './components/Drawer'
import Header from './components/Header'
import Home from './pages/Home'
import TrainingPage from './pages/TrainingPage'
import TrainingHistoryPage from './pages/TrainingsHistoryPage'
import CreateTrainingPage from './pages/CreateTrainingPage'
import DismissibleToast, { ToastMode } from './components/common/DismissibleToast'

const SessionCreatePage = lazy(() => import('./pages/SessionCreatePage'))

const [openToast, setOpenToast] = createSignal(false)
const [toastMessage, setToastMessage] = createSignal('')
const [toastMode, setToastMode] = createSignal(ToastMode.SUCCESS)

const showToast = (message: string, mode: ToastMode = ToastMode.SUCCESS) => {
  setToastMessage(message)
  setToastMode(mode)
  setOpenToast(true)
}

const App: Component = () => {
  return (
    <div>
      <Header />
      <Drawer />
      <DismissibleToast
        open={openToast()}
        onDismiss={() => setOpenToast(false)}
        mode={toastMode()}
        message={toastMessage()}
      />
      <div class="py-2">
        <Routes>
          <Route path="/" component={Home} />
          <Route path="/training/new" component={CreateTrainingPage} />
          <Route path="/training/:id" component={TrainingPage} />
          <Route path="/trainings" component={TrainingHistoryPage} />
          <Route path="/session/create" component={SessionCreatePage} />
        </Routes>
      </div>
    </div>
  )
}

export default App
export { showToast }
