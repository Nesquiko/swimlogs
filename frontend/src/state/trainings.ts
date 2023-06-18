import { createResource } from 'solid-js'
import {
  GetDetailsCurrWeekResponse,
  TrainingApi,
  TrainingDetail
} from '../generated'
import config from './api'

const trainingApi = new TrainingApi(config)

async function getTrainingsThisWeek(): Promise<GetDetailsCurrWeekResponse> {
  const result = trainingApi.getTrainingsDetailsCurrentWeek()
  const trainings = result.then((res) => res)
  return trainings
}
const [details, { mutate }] = createResource(getTrainingsThisWeek)
const useTrainingsDetails = () => [details]

function addTrainingDetail(td: TrainingDetail) {
  if (!isThisInThisWeek(td.date)) {
    return
  }
  const currentDetails = details()?.details ?? []
  const newDetails = [...currentDetails, td]
  newDetails.sort(trainingDetailCompare)
  mutate({ details: newDetails })
}

function isThisInThisWeek(date: Date): boolean {
  const todayObj = new Date()
  const todayDate = todayObj.getDate()
  const todayDay = (todayObj.getDay() - 1) % 7

  // get first date of week
  const firstDayOfWeek = new Date(todayObj.setDate(todayDate - todayDay))

  // get last date of week
  const lastDayOfWeek = new Date(firstDayOfWeek)
  lastDayOfWeek.setDate(lastDayOfWeek.getDate() + 6)

  // if date is equal or within the first and last dates of the week
  return date >= firstDayOfWeek && date <= lastDayOfWeek
}

function trainingDetailCompare(a: TrainingDetail, b: TrainingDetail): number {
  if (a.date < b.date) {
    return -1
  } else if (a.date > b.date) {
    return 1
  } else {
    const [aHours, aMinutes] = a.startTime.split(':').map(Number)
    const [bHours, bMinutes] = b.startTime.split(':').map(Number)

    if (aHours < bHours) {
      return -1
    } else if (aHours > bHours) {
      return 1
    } else {
      if (aMinutes < bMinutes) {
        return -1
      } else if (aMinutes > bMinutes) {
        return 1
      } else {
        return 0
      }
    }
  }
}

export { trainingApi, useTrainingsDetails, addTrainingDetail }
