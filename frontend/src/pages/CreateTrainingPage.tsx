import { useNavigate } from '@solidjs/router'
import {
  Component,
  createEffect,
  createResource,
  createSignal,
  Show
} from 'solid-js'
import { createStore } from 'solid-js/store'
import { Dynamic } from 'solid-js/web'
import { BlocksForm } from '../components/fragments/BlocksForm'
import { TrainingSessionForm } from '../components/fragments/TrainingSessionFormFragment'
import { CreateTrainingPreview } from '../components/fragments/TrainingPreviewFragment'
import { openToast, ToastType } from '../components/Toast'
import { NewTraining, ResponseError, StartingRuleType } from '../generated'
import { addTrainingDetail, trainingApi } from '../state/trainings'
import { sessionApi } from '../state/session'
import Spinner from '../components/Spinner'
import { CreateTrainingContextProvider } from '../components/context/CreateTrainingContextProvider'
import { NullDate, NullStartTime } from '../lib/consts'
import { useTransContext } from '@mbarzda/solid-i18next'

const CreateTrainingPage: Component = () => {
  const [training, setTraining] = createStore<NewTraining>({
    date: NullDate,
    startTime: NullStartTime,
    durationMin: 60,
    totalDistance: 100,
    blocks: [
      {
        num: 0,
        repeat: 1,
        name: '',
        totalDistance: 100,
        sets: [
          {
            num: 0,
            repeat: 1,
            distance: 100,
            what: '',
            startingRule: {
              type: StartingRuleType.None
            }
          }
        ]
      }
    ]
  })

  createEffect(() => {
    const totalDistance = training.blocks.reduce((acc, block) => {
      return acc + block.totalDistance
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
        console.error('error', e)
        openToast(
          t('training.creation.error', 'Error creating training'),
          ToastType.ERROR
        )
        navigate('/', { replace: true })
      })
  }

  const [currentComponent, setCurrentComponent] = createSignal(0)
  const comps = [TrainingSessionForm, BlocksForm, CreateTrainingPreview]

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
