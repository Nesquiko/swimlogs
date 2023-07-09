import { createEffect, For, JSX, on, onMount } from 'solid-js'

type MenuModalProps<O> = {
  opener: O
  items: { label: string; action: (o: O) => void }[]
  widthRem?: string
  header?: (o: O) => string
}

function MenuModal<O = {}>(props: MenuModalProps<O>): JSX.Element {
  let dialog: HTMLDialogElement
  const widthRem = props.widthRem ? props.widthRem + 'rem' : '11rem'

  createEffect(
    on(
      () => props.opener,
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
    <dialog style={{ width: widthRem }} ref={dialog!} class="rounded-lg">
      <p class="text-xl">{props.header?.(props.opener)}</p>
      <ul class="text-base text-black">
        <For each={props.items}>
          {(item) => {
            return (
              <li>
                <p
                  class="block cursor-pointer p-2 hover:bg-slate-200"
                  onClick={() => {
                    item.action(props.opener)
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
