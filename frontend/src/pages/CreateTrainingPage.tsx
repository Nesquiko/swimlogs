import { useNavigate } from '@solidjs/router'
import { Component, createEffect, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import { Dynamic } from 'solid-js/web'
import { NewTraining, ResponseError, StartType } from '../generated'
import { addTrainingDetail, trainingApi } from '../state/trainings'
import { NullDateTime } from '../lib/consts'
import { useTransContext } from '@mbarzda/solid-i18next'
import { openToast, ToastType } from '../components/Toast'
import TrainingSetsForm from './TrainingSetsForm'
import { TrainingSessionForm } from './TrainingSessionForm'
import { TrainingCreatePreviewPage } from './TrainingCreatePreviewPage'
import { CreateTrainingContextProvider } from '../components/CreateTrainingContextProvider'

const CreateTrainingPage: Component = () => {
  const [training, setTraining] = createStore<NewTraining>({
    start: NullDateTime,
    durationMin: 60,
    totalDistance: 100,
    sets: [
      {
        repeat: 1,
        setOrder: 0,
        distanceMeters: 100,
        totalDistance: 100,
        startType: StartType.None
      }
    ]
  })

  createEffect(() => {
    const totalDistance = training.sets.reduce((acc, set) => {
      return acc + set.totalDistance
    }, 0)
    setTraining('totalDistance', totalDistance)
  })

  const [t] = useTransContext()
  const navigate = useNavigate()
  async function createTraining() {
    const res = trainingApi.createTraining({ newTraining: training })
    await res
      .then((res) => {
        addTrainingDetail(res)
        openToast(t('training.created', 'Training created'), ToastType.SUCCESS)
        navigate('/', { replace: true })
      })
      .catch((e: ResponseError) => {
        // TODO check if error is 400
        console.error('error', e)
        openToast(
          t('training.creation.error', 'Error creating training'),
          ToastType.ERROR
        )
        navigate('/', { replace: true })
      })
  }

  const [currentComponent, setCurrentComponent] = createSignal(0)
  const comps = [
    TrainingSessionForm,
    TrainingSetsForm,
    TrainingCreatePreviewPage
  ]

  return (
    <div>
      <CreateTrainingContextProvider
        newTraining={training}
        currentComponentSignal={[currentComponent, setCurrentComponent]}
        sumbitTraining={createTraining}
      >
        <Dynamic component={comps[currentComponent()]} />
      </CreateTrainingContextProvider>
    </div>
  )
}

export default CreateTrainingPage
