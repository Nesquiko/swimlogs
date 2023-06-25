import { A } from '@solidjs/router'
import { Component } from 'solid-js'
import { openDrawer } from './Drawer'

const Topbar: Component = () => {
  return (
    <div class="w-full bg-sky-500 p-2">
      <i
        class="fa-solid fa-bars fa-xl inline-block w-1/5 cursor-pointer px-2 text-white"
        onClick={() => openDrawer()}
      ></i>
      <div class="inline-block w-3/5 text-center">
        <A href="/" class="cursor-pointer text-xl font-bold text-white">
          SwimLogs
        </A>
      </div>
    </div>
  )
}

export default Topbar
