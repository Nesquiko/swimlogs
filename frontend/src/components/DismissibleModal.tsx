import { Component, createEffect, onMount, Setter } from 'solid-js'

interface DismissibleModalProps {
  open: boolean
  setOpen: Setter<boolean>
  message: string
  confirmLabel: string
  cancelLabel: string

  onConfirm?: () => void
  onCancel?: () => void
}

const DismissibleModal: Component<DismissibleModalProps> = (props) => {
  let dialog: HTMLDialogElement

  const cancelDialog = () => {
    dialog.close()
    props.onCancel?.()
    props.setOpen(false)
  }

  const confirmDialog = () => {
    dialog.close()
    props.onConfirm?.()
    props.setOpen(false)
  }

  onMount(() => {
    dialog.addEventListener('click', (e) => {
      const dialogDimensions = dialog.getBoundingClientRect()
      if (
        e.clientX < dialogDimensions.left ||
        e.clientX > dialogDimensions.right ||
        e.clientY < dialogDimensions.top ||
        e.clientY > dialogDimensions.bottom
      ) {
        cancelDialog()
      }
    })
  })

  createEffect(() => {
    if (props.open) {
      dialog.showModal()
    }
  })

  return (
    <dialog ref={dialog!} class="max-w-xl w-11/12 rounded-lg shadow">
      <button
        type="button"
        class="absolute end-2.5 top-3 ms-auto inline-flex h-8 w-8 items-center justify-center rounded-lg bg-transparent text-sm text-gray-400 hover:bg-gray-200 hover:text-gray-900"
        onClick={cancelDialog}
      >
        <i class="fa-solid fa-xmark fa-xl"></i>
        <span class="sr-only">Close modal</span>
      </button>
      <div class="p-4 text-center md:p-5">
        <svg
          class="mx-auto mb-4 h-12 w-12 text-gray-400 dark:text-gray-200"
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
          onClick={confirmDialog}
        >
          {props.confirmLabel}
        </button>
        <button
          type="button"
          class="rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-base font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-900 focus:z-10 focus:outline-none focus:ring-4 focus:ring-gray-200 dark:border-gray-500 dark:bg-gray-700 dark:text-gray-300 dark:hover:bg-gray-600 dark:hover:text-white dark:focus:ring-gray-600"
          onClick={cancelDialog}
        >
          {props.cancelLabel}
        </button>
      </div>
    </dialog>
  )
}

export default DismissibleModal
