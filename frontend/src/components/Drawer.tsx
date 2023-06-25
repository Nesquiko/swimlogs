import { Trans } from '@mbarzda/solid-i18next'
import { A } from '@solidjs/router'
import { Component, createSignal } from 'solid-js'

const [open, setOpen] = createSignal(false)

const openDrawer = () => {
  setOpen(true)
}

const Drawer: Component = () => {
  return (
    <div>
      <div
        classList={{ 'bg-black/50': open(), 'pointer-events-none': !open() }}
        class="fixed left-0 top-0 h-full w-screen"
        onClick={() => {
          console.log('click')
        }}
      ></div>
      <div
        classList={{
          'translate-x-0': open(),
          '-translate-x-full': !open()
        }}
        class="fixed left-0 top-0 h-full w-full transform transition-transform duration-300 ease-in-out"
        /* onClick={() => setOpen(!open())} */
      >
        <div class="flex h-full w-1/2 flex-col justify-start space-y-8 bg-white pl-2 pt-4 md:w-1/4 lg:w-1/6">
          <A href="/training/create">
            <i class="fa-solid fa-person-swimming fa-2xl m-2 text-black"></i>
            <span class="font-bold text-black md:text-xl">
              <Trans key="create.training" />
            </span>
          </A>
          <A href="/session/create">
            <i class="fa-regular fa-clock fa-2xl m-2 text-black"></i>
            <span class="font-bold text-black md:text-xl">
              <Trans key="create.session" />
            </span>
          </A>
        </div>
      </div>
    </div>
  )
}

export { Drawer, openDrawer }
