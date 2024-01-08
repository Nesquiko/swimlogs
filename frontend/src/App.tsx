import { Routes, Route } from '@solidjs/router'
import { Component, lazy } from 'solid-js'
import { Drawer } from './components/Drawer'
import { Toast } from './components/Toast'
import Topbar from './components/Topbar'
import Home from './pages/Home'
import TrainingPage from './pages/TrainingPage'
import TrainingHistoryPage from './pages/TrainingsHistoryPage'

const TrainingCreatePage = lazy(() => import('./pages/TrainingCreatePage'))
const SessionCreatePage = lazy(() => import('./pages/SessionCreatePage'))

const App: Component = () => {
  return (
    <div>
      <Topbar />
      <Drawer />
      <div class="py-2">
        <Routes>
          <Route path="/" component={Home} />
          <Route path="/training/create" component={TrainingCreatePage} />
          <Route path="/training/:id" component={TrainingPage} />
          <Route path="/trainings" component={TrainingHistoryPage} />
          <Route path="/session/create" component={SessionCreatePage} />
        </Routes>
        <Toast />
      </div>
    </div>
  )
}

export default App
