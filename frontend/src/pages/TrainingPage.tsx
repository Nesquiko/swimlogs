import { useNavigate, useParams } from '@solidjs/router'
import { Component, createResource, Show } from 'solid-js'
import Spinner from '../components/Spinner'
import { openToast, ToastType } from '../components/Toast'
import TrainingPreview from '../components/TrainingPreview'
import { ResponseError, Training } from '../generated'
import { trainingApi } from '../state/trainings'

const TrainingPage: Component = () => {
  const params = useParams()
  const [training] = createResource(() => params.id, getTraining)

  const navigate = useNavigate()
  async function getTraining(id: string): Promise<Training> {
    const result = trainingApi
      .getTrainingById({ id })
      .then((res) => {
        return Promise.resolve(res)
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        let msg = 'Server error'
        if (e.response?.status === 404) {
          msg = 'Training not found'
        }
        openToast(msg, ToastType.ERROR)
        navigate('/', { replace: true })
        return Promise.reject(e)
      })
    return result
  }

  return (
    <Show when={!training.loading} fallback={<Spinner remSize={8} />}>
      <TrainingPreview training={training()} />
    </Show>
  )
}

export default TrainingPage
