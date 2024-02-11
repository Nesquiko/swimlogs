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

/** Monday = 1, Tuesday = 2, ..., Sunday = 7 */
export function dateForDayOfWeek(targetDayIndex: number): Date {
  const today = new Date()

  const curr = new Date()
  const first = curr.getDate() - (6 - curr.getDay())
  const target = first + targetDayIndex - 1

  const date = new Date(today.getFullYear(), today.getMonth(), target)

  return new Date(date)
}

export function isThisInThisWeek(date: Date): boolean {
  const curr = new Date()
  const first = curr.getDate() - (6 - curr.getDay())
  const last = first + 6

  return first <= date.getDate() && date.getDate() <= last
}
