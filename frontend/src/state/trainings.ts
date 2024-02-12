import { createResource } from 'solid-js'
import {
  BASE_PATH,
  Configuration,
  TrainingDetail,
  TrainingDetailsCurrentWeekResponse,
  TrainingsApi,
} from 'swimlogs-api'
import { isThisInThisWeek } from '../lib/datetime'

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
  newDetails.sort(trainingDetailCompare)
  mutate({ details: newDetails })
}

function removeFromTrainingsDetails(id: string) {
  const currentDetails = detailsThisWeek()?.details ?? []
  const newDetails = currentDetails.filter((td) => td.id !== id)
  mutate({ details: newDetails })
}

function updateTrainintDetails(td: TrainingDetail) {
  if (!isThisInThisWeek(td.start)) {
    removeFromTrainingsDetails(td.id)
    return
  }

  const details = detailsThisWeek()?.details ?? []
  const tdIdx = Math.max(
    details.findIndex((detail) => detail.id === td.id),
    0
  )
  details[tdIdx] = td
  details.sort(trainingDetailCompare)
  mutate({ details: details })
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
  updateTrainintDetails,
}
