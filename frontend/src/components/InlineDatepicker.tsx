import { useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, For } from 'solid-js'
import { capitalize } from '../lib/str'

function locale() {
  if (navigator.languages != undefined) return navigator.languages[0]
  return navigator.language
}

const DAY_NAME_KEYS = [
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
  'sunday',
]

type CalendarDay = {
  day: number
  month: 'next' | 'current' | 'previous'
}

const MONDAY = 1
const SUNDAY = 0

interface InlineDatepickerProps {
  initialDate?: Date
  onChange: (date: Date) => void
}

const InlineDatepicker: Component<InlineDatepickerProps> = (props) => {
  const [t] = useTransContext()
  const [date, _setDate] = createSignal(
	props.initialDate || new Date()
  )
  const lastDay = () =>
    new Date(date().getUTCFullYear(), date().getUTCMonth() + 1, 1)

  const setDate = (newDate: Date) => {
	_setDate(newDate)
	props.onChange(newDate)
  }

  const days = () => {
    const days = new Array<CalendarDay>()
    const firstDay = new Date(date().getFullYear(), date().getMonth(), 1)

    if (firstDay.getDay() !== MONDAY) {
      let fromPrevious = new Date(date().getFullYear(), date().getMonth(), 0)

      while (fromPrevious.getDay() !== SUNDAY) {
        days.push({
          day: fromPrevious.getDate(),
          month: 'previous',
        })

        fromPrevious = new Date(
          fromPrevious.getFullYear(),
          fromPrevious.getMonth(),
          fromPrevious.getDate() - 1
        )
      }
    }
    days.reverse()

    for (let i = 0; i < lastDay().getUTCDate(); i++) {
      days.push({
        day: i + 1,
        month: 'current',
      })
    }

    let nextDay = new Date(date().getFullYear(), date().getMonth() + 1, 1)
    while (nextDay.getDay() !== MONDAY) {
      days.push({
        day: nextDay.getDate(),
        month: 'next',
      })

      nextDay = new Date(
        nextDay.getFullYear(),
        nextDay.getMonth(),
        nextDay.getDate() + 1
      )
    }

    return days
  }

  const subtractMonth = () => {
    setDate(new Date(date().setMonth(date().getMonth() - 1)))
  }

  const addMonth = () => {
    setDate(new Date(date().setMonth(date().getMonth() + 1)))
  }

  const dayButton = (day: CalendarDay) => {
    return (
      <button
        classList={{
          'text-gray-500': day.month !== 'current',
          '!bg-gray-400':
            date().getDate() === day.day && day.month === 'current',
        }}
        class="h-10 rounded-lg text-center align-middle leading-10 hover:bg-gray-200"
        onClick={() => {
          if (day.month === 'previous') {
            const newDate = new Date()
			newDate.setMonth(date().getMonth() - 1)
            newDate.setDate(day.day)
            setDate(newDate)
          } else if (day.month === 'next') {
            const newDate = new Date()
			newDate.setMonth(date().getMonth() + 1)
            newDate.setDate(day.day)
            setDate(newDate)
          } else {
            setDate(new Date(date().setDate(day.day)))
          }
        }}
      >
        <p>{day.day}</p>
      </button>
    )
  }

  return (
    <div class="w-full space-y-2 p-4 md:w-96">
      <div class="flex justify-between">
        <button onClick={subtractMonth}>
          <i class="fas fa-chevron-left fa-xl"></i>
        </button>

        <p class="text-lg font-bold">
          {t(
            date().toLocaleDateString(locale(), { month: 'long' }).toLowerCase()
          )}
          {', '}
          {date().getFullYear()}
        </p>

        <button onClick={addMonth}>
          <i class="fas fa-chevron-right fa-xl"></i>
        </button>
      </div>

      <div class="grid grid-cols-7">
        {DAY_NAME_KEYS.map((dayName) => (
          <div class="text-center">
            <p class="text-gray-500">
              {capitalize((t(dayName) || dayName).slice(0, 2))}
            </p>
          </div>
        ))}
      </div>

      <div class="grid grid-cols-7">
        <For each={days()}>{dayButton}</For>
      </div>
    </div>
  )
}

export default InlineDatepicker
