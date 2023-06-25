import { Routes, Route } from '@solidjs/router'
import { Component, lazy } from 'solid-js'
import { Drawer } from './components/Drawer'
import { Toast } from './components/Toast'
import Topbar from './components/Topbar'
import Home from './pages/Home'
import TrainingPage from './pages/TrainingPage'

const CreateTrainingPage = lazy(() => import('./pages/CreateTrainingPage'))
const SessionCreatePage = lazy(() => import('./pages/SessionCreatePage'))

const App: Component = () => {
  return (
    <div>
      <Topbar />
      <Drawer />
      <Routes>
        <Route path="/" component={Home} />
        <Route path="/training/create" component={CreateTrainingPage} />
        <Route path="/training/:id" component={TrainingPage} />
        <Route path="/session/create" component={SessionCreatePage} />
      </Routes>
      <Toast />
    </div>
  )
}

export default App
