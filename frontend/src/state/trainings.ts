import { createResource } from 'solid-js'
import {
  BASE_PATH,
  Configuration,
  TrainingDetail,
  TrainingDetailsCurrentWeekResponse,
  TrainingsApi,
} from 'swimlogs-api'

const config = new Configuration({
  basePath: import.meta.env.DEV ? 'http://localhost:42069' : BASE_PATH,
})

const trainingApi = new TrainingsApi(config)

async function getTrainingsThisWeek(): Promise<TrainingDetailsCurrentWeekResponse> {
  return trainingApi.trainingDetailsCurrentWeek().catch((_err) => {
    throw Promise.reject(new Error('Error fetching trainings details'))
  })
}
const [detailsThisWeek, { mutate }] = createResource(getTrainingsThisWeek)
const useTrainingsDetailsThisWeek = () => [detailsThisWeek]

function addTrainingDetail(td: TrainingDetail) {
  if (!isThisInThisWeek(td.start)) {
    return
  }
  const currentDetails = detailsThisWeek()?.details ?? []
  const newDetails = [...currentDetails, td]
  const newDistance = newDetails.reduce(
    (acc, curr) => acc + curr.totalDistance,
    0
  )
  newDetails.sort(trainingDetailCompare)
  mutate({ details: newDetails, distance: newDistance })
}

function removeFromTrainingsDetails(id: string) {
  const currentDetails = detailsThisWeek()?.details ?? []
  const newDetails = currentDetails.filter((td) => td.id !== id)
  const newDistance = newDetails.reduce(
    (acc, curr) => acc + curr.totalDistance,
    0
  )
  mutate({ details: newDetails, distance: newDistance })
}

function isThisInThisWeek(date: Date): boolean {
  const todayObj = new Date()
  const todayDate = todayObj.getDate()
  const todayDay = (todayObj.getDay() - 1) % 7

  // get first date of week
  const firstDayOfWeek = new Date(todayObj.setDate(todayDate - todayDay))
  firstDayOfWeek.setHours(0, 0, 0, 0)

  // get last date of week
  const lastDayOfWeek = new Date(firstDayOfWeek)
  lastDayOfWeek.setDate(lastDayOfWeek.getDate() + 6)
  lastDayOfWeek.setHours(23, 59, 59, 999)

  // if date is equal or within the first and last dates of the week
  return date >= firstDayOfWeek && date <= lastDayOfWeek
}

function trainingDetailCompare(a: TrainingDetail, b: TrainingDetail): number {
  if (a.start < b.start) {
    return -1
  } else if (a.start > b.start) {
    return 1
  }
  return 0
}

export {
  trainingApi,
  useTrainingsDetailsThisWeek,
  addTrainingDetail,
  removeFromTrainingsDetails,
}
