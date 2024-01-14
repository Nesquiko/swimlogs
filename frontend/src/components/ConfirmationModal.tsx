import { Modal } from 'flowbite'
import { Component, onMount } from 'solid-js'
import { randomId } from '../lib/str'

interface ConfirmationModalProps {
  icon: string
  message: string

  confirmLabel: string
  cancelLabel: string

  onConfirm: () => void
  onCancel: () => void
}

const ConfirmationModal: Component<ConfirmationModalProps> = (props) => {
  const id = randomId()
  let target: HTMLDivElement
  let modal: Modal

  onMount(() => {
    modal = new Modal(target!, {}, { id })
  })

  return (
    <>
      <button
        class="rounded-lg p-1 text-sky-900"
        type="button"
        onClick={() => modal.show()}
      >
        <i class={`${props.icon} fa-solid fa-xl cursor-pointer`}></i>
      </button>

      <div
        id={id}
        ref={target!}
        tabindex="-1"
        class="fixed left-0 right-0 top-0 z-50 hidden h-[calc(100%-1rem)] max-h-full w-full items-center justify-center overflow-y-auto overflow-x-hidden md:inset-0"
      >
        <div class="relative max-h-full w-full max-w-md p-4">
          <div class="relative rounded-lg bg-white shadow dark:bg-gray-700">
            <button
              type="button"
              class="absolute end-2.5 top-3 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-transparent text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-white"
              onClick={() => modal.hide()}
            >
              <i class="fa-solid fa-xmark fa-xl"></i>
              <span class="sr-only">Close modal</span>
            </button>
            <div class="p-4 text-center md:p-5">
              <svg
                class="mx-auto mb-4 h-12 w-12 text-gray-400 dark:text-gray-200"
                aria-hidden="true"
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 20 20"
              >
                <path
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M10 11V6m0 8h.01M19 10a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
                />
              </svg>
              <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
                {props.message}
              </h3>
              <button
                type="button"
                class="me-2 inline-flex items-center rounded-lg bg-red-500 px-5 py-2.5 text-center text-base font-medium text-white hover:bg-red-800 focus:outline-none focus:ring-4 focus:ring-red-300 dark:focus:ring-red-800"
                onClick={() => {
                  props.onConfirm()
                  modal.hide()
                }}
              >
                {props.confirmLabel}
              </button>
              <button
                type="button"
                class="rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-base font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-900 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:border-gray-500 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white dark:focus:ring-gray-600"
                onClick={() => {
                  props.onCancel()
                  modal.hide()
                }}
              >
                {props.cancelLabel}
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  )
}

export default ConfirmationModal
