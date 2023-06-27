import { Trans } from '@mbarzda/solid-i18next'
import { useLocation, useNavigate } from '@solidjs/router'
import { Component, createSignal } from 'solid-js'
import swimSvg from '../assets/swim.svg'
import clockSvg from '../assets/clock.svg'
import homeSvg from '../assets/home.svg'

const [open, setOpen] = createSignal(false)

const openDrawer = () => {
  setOpen(true)
}

interface ItemProps {
  pathname: string
  labelKey: string
  imgSrc: string
}

const Drawer: Component = () => {
  const location = useLocation()
  const navigate = useNavigate()

  const Item: Component<ItemProps> = (props) => {
    return (
      <div
        classList={{ 'bg-sky-500/20': location.pathname === props.pathname }}
        class="cursor-pointer p-2"
        onClick={() => {
          setOpen(false)
          navigate(props.pathname)
        }}
      >
        <img class="inline h-8 w-8 md:h-10 md:w-10" src={props.imgSrc} />
        <span class="pl-2 text-lg font-bold text-black">
          <Trans key={props.labelKey} />
        </span>
      </div>
    )
  }

  return (
    <div>
      <div
        classList={{ 'bg-black/50': open(), 'pointer-events-none': !open() }}
        class="fixed left-0 top-0 z-10 h-full w-full"
        onClick={() => setOpen(!open())}
      ></div>
      <div
        classList={{
          'translate-x-0': open(),
          '-translate-x-full': !open()
        }}
        class="fixed left-0 top-0 z-20 h-full w-1/2 transform transition-transform duration-300 ease-in-out sm:w-1/3 md:w-1/4 lg:w-1/5 xl:w-1/6"
      >
        <div class="h-full flex-col justify-start bg-white pt-4">
          <Item pathname="/" labelKey="home" imgSrc={homeSvg} />
          <Item
            pathname="/training/create"
            labelKey="create.training"
            imgSrc={swimSvg}
          />
          <Item
            pathname="/session/create"
            labelKey="create.session"
            imgSrc={clockSvg}
          />
        </div>
      </div>
    </div>
  )
}

export { Drawer, openDrawer }
