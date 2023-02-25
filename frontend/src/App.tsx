import { Routes, Route } from 'solid-app-router'
import type { Component } from 'solid-js'
import Topbar from './component/Topbar'
import Home from './pages/Home'

const App: Component = () => {
  return (
    <div>
      <Topbar />
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </div>
  )
}

export default App
