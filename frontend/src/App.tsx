import 'flowbite'
import { Routes, Route } from '@solidjs/router'
import { Component, lazy } from 'solid-js'
import { Drawer } from './components/Drawer'
import Header from './components/Header'
import { Toast } from './components/Toast'
import Home from './pages/Home'
import TrainingPage from './pages/TrainingPage'
import TrainingHistoryPage from './pages/TrainingsHistoryPage'
import CreateTrainingPage from './pages/CreateTrainingPage'

const TrainingCreatePage = lazy(() => import('./pages/TrainingCreatePage'))
const SessionCreatePage = lazy(() => import('./pages/SessionCreatePage'))

const App: Component = () => {
  return (
    <div>
      <Header />
      <Drawer />
      <div class="py-2">
        <Routes>
          <Route path="/" component={Home} />
          <Route path="/training/create" component={TrainingCreatePage} />
          <Route path="/training/new" component={CreateTrainingPage} />
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
