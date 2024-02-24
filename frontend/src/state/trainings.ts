import { RouteLoadFunc } from '@solidjs/router'
import { createResource } from 'solid-js'
import {
  BASE_PATH,
  Configuration,
  ResponseError,
  Training,
  TrainingDetail,
  TrainingDetailsCurrentWeekResponse,
  TrainingsApi,
} from 'swimlogs-api'
import { isThisInThisWeek } from '../lib/datetime'

const config = new Configuration({
  basePath: import.meta.env.DEV ? 'http://localhost:42069' : BASE_PATH,
})

const trainingApi = new TrainingsApi(config)

type TrainingCacheEntry = { id: string; t: Promise<Training> }
let trainingCache: { entry: TrainingCacheEntry | undefined } = {
  entry: undefined,
}

export const loadTrainingById: RouteLoadFunc<Promise<Training>> = ({
  params,
}): Promise<Training> => {
  if (trainingCache.entry && trainingCache.entry.id === params.id) {
    return trainingCache.entry.t
  }

  const t = trainingApi
    .training({ id: params.id })
    .catch((e: ResponseError) => {
      trainingCache.entry = undefined
      throw e
    })
  trainingCache.entry = { id: params.id, t }
  return t
}

export const deleteTrainingById = async (id: string): Promise<void> => {
  return await trainingApi
    .deleteTraining({ id })
    .then(() => removeFromTrainingsDetails(id))
    .catch(() => {
      removeFromTrainingsDetails(id)
    })
    .finally(() => (trainingCache.entry = undefined))
}

export const updateTrainingById = async (
  id: string,
  training: Training
): Promise<TrainingDetail> => {
  return trainingApi
    .editTraining({ id, training })
    .then((res) => {
      updateTrainintDetails(res)
      trainingCache.entry = undefined
      return res
    })
    .catch((e: ResponseError) => {
      trainingCache.entry = undefined
      throw e
    })
}

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
