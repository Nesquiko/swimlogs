import { expect, test } from 'vitest'
import { isThisInWeek } from '../src/lib/datetime'

test('per day matrix, all in given week', () => {
  const date = new Date(2024, 1, 12, 12) // Monday

  for (const d of [0, 1, 2, 3, 4, 5, 6]) {
    const referenceDate = addDays(date, d)
    for (const d of [0, 1, 2, 3, 4, 5, 6]) {
      const day = addDays(date, d)
      expect(isThisInWeek(referenceDate, day)).toBeTruthy()
    }
  }
})

function addDays(date: Date, days: number): Date {
  var result = new Date(date)
  result.setDate(result.getDate() + days)
  return result
}
