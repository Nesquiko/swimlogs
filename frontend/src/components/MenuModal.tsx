import { Component, createEffect, For, on, onMount } from 'solid-js'

type MenuItem = {
  label: string
  action: () => void
}

type MenuModalProps = {
  items: MenuItem[]
  // indicates that the dialog should be opened, since the dialog can close itself,
  // we only need to indicate when to open it
  open: {}
}

const MenuModal: Component<MenuModalProps> = (props) => {
  let dialog: HTMLDialogElement

  createEffect(
    on(
      () => props.open,
      () => {
        dialog.inert = true // disables auto focus when dialog is opened
        dialog.showModal()
        dialog.inert = false
      },
      { defer: true }
    )
  )

  onMount(() => {
    dialog.addEventListener('click', (e) => {
      const dialogDimensions = dialog.getBoundingClientRect()
      if (
        e.clientX < dialogDimensions.left ||
        e.clientX > dialogDimensions.right ||
        e.clientY < dialogDimensions.top ||
        e.clientY > dialogDimensions.bottom
      ) {
        dialog.close()
      }
    })
  })

  return (
    <dialog ref={dialog!} class="w-44 rounded-lg">
      <ul class="text-md text-black">
        <For each={props.items}>
          {(item) => {
            return (
              <li>
                <p
                  class="block cursor-pointer p-2 hover:bg-slate-200"
                  onClick={() => {
                    item.action()
                    dialog.close()
                  }}
                >
                  {item.label}
                </p>
              </li>
            )
          }}
        </For>
      </ul>
    </dialog>
  )
}

export default MenuModal
