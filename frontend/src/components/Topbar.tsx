import { A } from '@solidjs/router'
import { Component } from 'solid-js'
import { openDrawer } from './Drawer'

const Topbar: Component = () => {
  return (
    <div class="w-full bg-sky-500 p-2">
      <div class="inline-block w-1/5 align-middle">
        <img
          class="inline-block cursor-pointer"
          src="/src/assets/menu.svg"
          width={40}
          height={40}
          onClick={() => openDrawer()}
        />
      </div>
      <div class="inline-block w-3/5 text-center align-middle">
        <A href="/" class="cursor-pointer text-xl font-bold text-white">
          SwimLogs
        </A>
      </div>
    </div>
  )
}

export default Topbar
