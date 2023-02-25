import { Routes, Route } from 'solid-app-router'
import type { Component } from 'solid-js'
import Topbar from './component/Topbar'
import Home from './pages/Home'
import CreateTraining from './pages/CreateTraining'

const App: Component = () => {
  return (
    <div>
      <Topbar />
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/training/create" element={<CreateTraining />} />
      </Routes>
    </div>
  )
}

export default App
