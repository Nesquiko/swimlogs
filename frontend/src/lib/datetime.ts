export const DayEnum = {
  Monday: 'Monday',
  Tuesday: 'Tuesday',
  Wednesday: 'Wednesday',
  Thursday: 'Thursday',
  Friday: 'Friday',
  Saturday: 'Saturday',
  Sunday: 'Sunday',
} as const
export type DayEnum = (typeof DayEnum)[keyof typeof DayEnum]

export function formatDate(date: Date | undefined): string {
  if (!date) {
    return ''
  }
  const day = String(date.getDate()).padStart(2, '0')
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const year = String(date.getFullYear())

  return `${day}.${month}.${year}`
}

export function formatTime(date: Date | undefined): string {
  if (!date) {
    return ''
  }
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')

  return `${hours}:${minutes}`
}

export function datesThisWeek(): Date[] {
  const curr = new Date()
  const daysInWeek: Date[] = []
  const dayOfWeek = curr.getDay()

  const firstDay = new Date(curr)
  firstDay.setDate(curr.getDate() - (dayOfWeek === 0 ? 6 : dayOfWeek - 1))

  for (let i = 0; i < 7; i++) {
    const currentDate = new Date(firstDay)
    currentDate.setDate(firstDay.getDate() + i)
    daysInWeek.push(currentDate)
  }

  return daysInWeek
}

export function isThisInThisWeek(date: Date): boolean {
  return isThisInWeek(new Date(), date)
}

const SUNDAY = 0

export function isThisInWeek(referenceDate: Date, date: Date): boolean {
  let first: number
  if (referenceDate.getDay() === SUNDAY) {
    first = referenceDate.getDate() - 6
  } else {
    first = referenceDate.getDate() + (-referenceDate.getDay() + 1)
  }

  const last = first + 6

  return first <= date.getDate() && date.getDate() <= last
}

export function locale() {
  if (navigator.languages != undefined) return navigator.languages[0]
  return navigator.language
}

export function minutesToHoursAndMintes(minutes: number): string {
  const hours = Math.floor(minutes / 60)
  const remainingMinutes = minutes % 60

  let result = ''

  if (hours > 0) {
    result += `${hours}h `
  }

  if (remainingMinutes > 0 || result === '') {
    result += `${remainingMinutes}min`
  }

  return result
}
