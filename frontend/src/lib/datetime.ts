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

export function dateForDayOfWeek(targetDayIndex: number): Date {
  const today = new Date()
  const todayIndex = today.getDay()
  const daysUntilTarget = targetDayIndex - todayIndex
  today.setDate(today.getDate() + daysUntilTarget)
  return new Date(today)
}
