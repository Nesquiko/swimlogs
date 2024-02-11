import { Dismiss } from 'flowbite'
import { Component, createEffect, onMount, Show } from 'solid-js'
import { randomId } from '../../lib/str'

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

const TRANSITION_DURATION = 300

// TODO the dismiss doesn't work... remake it so that I can spam the toas,
//  when the toast is triggered 2 in row, the second time it isn't shown...
const DismissibleToast: Component<DismissibleToastProps> = (props) => {
  const id = randomId()

  let target: HTMLDivElement
  let dismissEl: Dismiss

  onMount(() => {
    dismissEl = new Dismiss(
      target,
      undefined,
      { duration: TRANSITION_DURATION },
      { id }
    )
  })

  const onDismiss = () => {
    dismissEl.hide()
    setTimeout(props.onDismiss, TRANSITION_DURATION)
  }

  createEffect(() => {
    if (props.open) {
      dismissEl.removeInstance
      setTimeout(onDismiss, 3000)
    }
  })

  return (
    <div role="alert" ref={target!}>
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
