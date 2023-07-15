import { Trans } from '@mbarzda/solid-i18next'
import { createEffect, createSignal, For, JSX, on, onMount } from 'solid-js'
import { Day, Session } from '../generated'
import { formatDate } from '../lib/datetime'

type DatePickerModalProps = {
  opener: { s: Session }
  onSubmit: (date: Date, s: Session) => void
  onCancel?: () => void
}

function DatePickerModal(props: DatePickerModalProps): JSX.Element {
  let dialog: HTMLDialogElement
  const [dates, setDates] = createSignal<Date[]>([])

  createEffect(
    on(
      () => props.opener,
      () => {
        setDates(getDates(props.opener.s))
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

  const getDates = (s: Session): Date[] => {
    const dayNames = Array.from(Object.keys(Day))
    const inputDayIndex = dayNames.indexOf(s.day)
    const today = new Date()
    const todayIndex = today.getDay() - 1
    const dayDifference = (inputDayIndex - todayIndex + 7) % 7
    const futureDates: Date[] = []
    // Start from today and find the next four dates on the input day
    for (let i = 0; i < 4; i++) {
      const futureDate = new Date(
        today.getTime() + (dayDifference + 7 * i) * 24 * 60 * 60 * 1000
      )
      futureDate.setHours(0, 0, 0, 0)
      futureDates.push(futureDate)
    }
    return futureDates
  }

  return (
    <dialog ref={dialog!} class="rounded-lg">
      <p class="text-xl">
        <Trans key="select.date" />
      </p>

      <select
        id="date"
        class="my-2 rounded-md border border-solid bg-white px-4 py-2 text-xl focus:border-sky-500 focus:outline-none focus:ring"
        onChange={(e) => {
          const date = new Date(dates()[parseInt(e.target.value)])
          const hours = parseInt(props.opener.s.startTime.slice(0, 2))
          const minutes = parseInt(props.opener.s.startTime.slice(3))
          date.setHours(hours, minutes, 0, 0)
          props.onSubmit(date, props.opener.s)
          dialog.close()
        }}
      >
        <option disabled selected>
          <Trans key="select.date" />
        </option>
        <For each={dates()}>
          {(d, i) => <option value={i()}>{formatDate(d)}</option>}
        </For>
      </select>
    </dialog>
  )
}

export default DatePickerModal
