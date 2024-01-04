import { Component, createSignal } from 'solid-js'

enum ToastType {
  INFO,
  ERROR,
  SUCCESS,
}

const [open, setOpen] = createSignal(false)
const [move, setMove] = createSignal(false)
const [type, setType] = createSignal(ToastType.INFO)
const [text, setText] = createSignal('Toast')

const openToast = (
  text: string,
  type: ToastType = ToastType.INFO,
  timeout: number = 2000
) => {
  if (open()) return
  setType(type)
  setText(text)
  setOpen(true)
  setMove(true)
  setTimeout(() => setMove(false), timeout)
  setTimeout(() => setOpen(false), timeout + 150)
}

const Toast: Component = () => {
  return (
    <div
      id="toast"
      classList={{
        'translate-y-full': !move(),
        '-translate-y-full': move(),
        visible: open(),
        invisible: !open(),
        'h-0': !open(),
      }}
      class="fixed bottom-0 z-10 mx-auto flex w-full transform justify-center transition-transform duration-300 ease-in-out"
    >
      <div
        classList={{
          'bg-sky-500 border-sky-900': type() === ToastType.INFO,
          'bg-red-500 border-red-900': type() === ToastType.ERROR,
          'bg-green-500': type() === ToastType.SUCCESS,
        }}
        class="w-11/12 rounded-lg border p-2 text-xl font-bold text-white shadow"
      >
        {text()}
      </div>
    </div>
  )
}

export { Toast, ToastType, openToast }
