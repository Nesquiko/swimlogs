import { useTransContext } from '@mbarzda/solid-i18next'
import { useNavigate, useParams } from '@solidjs/router'
import { Component, createResource, Show } from 'solid-js'
import { showToast } from '../App'
import { ToastMode } from '../components/common/DismissibleToast'
import Spinner from '../components/Spinner'
import { ResponseError, Training } from '../generated'
import { removeFromTrainingsDetails, trainingApi } from '../state/trainings'
import TrainingPreviewPage from './TrainingPreviewPage'

const TrainingPage: Component = () => {
  const params = useParams()
  const [t] = useTransContext()
  const [training] = createResource(() => params.id, getTraining)

  const navigate = useNavigate()
  async function getTraining(id: string): Promise<Training> {
    return trainingApi
      .getTrainingById({ id })
      .then((res) => res)
      .catch((e: ResponseError) => {
        console.error('error', e)
        let msg = t('server.error')
        if (e.response?.status === 404) {
          msg = t('training.not.found')
        }
        showToast(msg, ToastMode.ERROR)
        navigate('/', { replace: true })
        return Promise.reject(e)
      })
  }

  async function deleteTraining(id: string) {
    trainingApi
      .deleteTraining({ id: id })
      .then(() => removeFromTrainingsDetails(id))
      .catch((e: ResponseError) => {
        console.error('error', e)
        let msg = t('server.error')
        if (e.response?.status === 404) {
          msg = t('training.not.found')
          removeFromTrainingsDetails(id)
        }
        showToast(msg, ToastMode.ERROR)
      })
      .finally(() => {
        navigate('/', { replace: true })
      })
  }

  return (
    <Show when={!training.error}>
      <Show when={!training.loading} fallback={<Spinner remSize={8} />}>
        <TrainingPreviewPage
          training={training()!}
          showDeleteTraining={true}
          onDeleteTraining={() => deleteTraining(params.id)}
        />
      </Show>
    </Show>
  )
}

export default TrainingPage
