import { Component, createEffect, createSignal, Show } from 'solid-js'

enum ToastMode {
  SUCCESS = 'SUCCESS',
  ERROR = 'ERROR',
}

const ToastModeOptions = {
  [ToastMode.SUCCESS]: { icon: 'fa-check ', color: 'text-green-500' },
  [ToastMode.ERROR]: { icon: 'fa-exclamation', color: 'text-red-500' },
} as const

interface DismissibleToastProps {
  open: boolean
  mode: ToastMode
  onDismiss: () => void
  message: string
}

const TRANSITION_DURATION = 1000
const AUTOMATIC_CLOSE_DURATION = 3000

const DismissibleToast: Component<DismissibleToastProps> = (props) => {
  let closeTimeout: NodeJS.Timeout
  const [easeOut, setEaseOut] = createSignal(false)

  const onDismiss = () => {
    clearTimeout(closeTimeout)
    setEaseOut(true)
    setTimeout(() => {
      props.onDismiss()
      setEaseOut(false)
    }, TRANSITION_DURATION)
  }

  createEffect(() => {
    if (props.open) {
      closeTimeout = setTimeout(onDismiss, AUTOMATIC_CLOSE_DURATION)
    }
  })

  return (
    <div
      role="alert"
      classList={{
        'opacity-100': !easeOut(),
        'opacity-0': easeOut(),
      }}
      class="transition-opacity ease-in duration-300"
    >
      <Show when={props.open}>
        <div class="fixed right-4 top-4 flex w-full max-w-xs items-center rounded-lg bg-white p-4 text-gray-800 shadow shadow-gray-400">
          <i
            class={`fa-solid fa-xl ${ToastModeOptions[props.mode].color} ${
              ToastModeOptions[props.mode].icon
            }`}
          ></i>
          <div class="ms-3 text-base font-medium">{props.message}</div>
          <button
            type="button"
            class="-mx-1.5 -my-1.5 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg p-1.5"
            aria-label="Close"
            onClick={onDismiss}
          >
            <i class="fa-solid fa-xmark fa-xl"></i>
          </button>
        </div>
      </Show>
    </div>
  )
}

export default DismissibleToast
export { ToastMode }
