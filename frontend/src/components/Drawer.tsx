import { useTransContext } from '@mbarzda/solid-i18next'
import { useLocation, useNavigate } from '@solidjs/router'
import { type Component, createSignal } from 'solid-js'

const [open, setOpen] = createSignal(false)

const openDrawer = (): void => {
  setOpen(true)
}

interface ItemProps {
  path: string
  label: string
  icon: string
}

const Drawer: Component = () => {
  const [t] = useTransContext()
  const location = useLocation()
  const navigate = useNavigate()

  const Item: Component<ItemProps> = (props) => {
    return (
      <div
        classList={{ 'bg-sky-500/20': location.pathname === props.path }}
        class="cursor-pointer p-2"
        onClick={() => {
          setOpen(false)
          navigate(props.path)
        }}
      >
        <i class={`fa-solid fa-xl ${props.icon} `}></i>
        <span class="pl-2 text-lg font-bold text-black">{props.label}</span>
      </div>
    )
  }

  return (
    <div id="drawer">
      <div
        classList={{ 'bg-black/50': open(), 'pointer-events-none': !open() }}
        class="fixed left-0 top-0 z-10 h-full w-full"
        onClick={() => setOpen(!open())}
      ></div>
      <div
        classList={{
          'translate-x-0': open(),
          '-translate-x-full': !open(),
        }}
        class="fixed left-0 top-0 z-20 h-full w-1/2 transform transition-transform duration-300 ease-in-out sm:w-1/3 md:w-1/4 lg:w-1/5 xl:w-1/6"
      >
        <div class="h-full bg-white pt-4">
          <Item path="/" label={t('home')} icon="fa-house" />
          <Item
            path="/trainings"
            label={t('trainings.history')}
            icon="fa-clock-rotate-left"
          />
        </div>
      </div>
    </div>
  )
}

export { Drawer, openDrawer }
