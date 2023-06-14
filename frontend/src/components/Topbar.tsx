import { A } from '@solidjs/router'
import { Component } from 'solid-js'

const Topbar: Component = () => {
  return (
    <div class="flex w-full items-center justify-center bg-sky-500 p-2">
      <A href="/">
        <span class="text-xl font-bold text-white">SwimLogs</span>
      </A>
    </div>
  )
}

export default Topbar
