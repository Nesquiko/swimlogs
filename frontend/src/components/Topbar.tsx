import { A } from '@solidjs/router'
import { Component } from 'solid-js'
import { openDrawer } from './Drawer'
import menuSvg from '../assets/menu.svg'

const Topbar: Component = () => {
  return (
    <div id="topbar" class="w-full bg-sky-500 p-2">
      <div class="inline-block w-1/5 align-middle">
        <i
          class="fa-solid fa-bars fa-2xl cursor-pointer"
          style="color: #ffffff;"
          onClick={openDrawer}
        ></i>
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
