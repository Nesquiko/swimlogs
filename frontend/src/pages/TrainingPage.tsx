import { useTransContext } from '@mbarzda/solid-i18next'
import { useNavigate, useParams } from '@solidjs/router'
import { type Component, createResource, createSignal, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import { showToast } from '../App'
import { ToastMode } from '../components/common/DismissibleToast'
import ConfirmationModal from '../components/ConfirmationModal'
import Spinner from '../components/Spinner'
import { ResponseError, Training } from 'swimlogs-api'
import {
  removeFromTrainingsDetails,
  trainingApi,
  updateTrainintDetails,
} from '../state/trainings'
import EditTrainingPage from './EditTrainingPage'
import TrainingPreviewPage from './TrainingPreviewPage'

const TrainingPage: Component = () => {
  const params = useParams()
  const [t] = useTransContext()
  const [training, { mutate }] = createResource(() => params.id, getTraining)
  const [editTraining, setEditTraining] = createSignal(false)

  const navigate = useNavigate()
  async function getTraining(id: string): Promise<Training> {
    return await trainingApi
      .training({ id })
      .then((res) => res)
      .catch((e: ResponseError) => {
        console.error('error', e)
        let msg = t('server.error')
        if (e.response?.status === 404) {
          msg = t('training.not.found')
        }
        showToast(msg, ToastMode.ERROR)
        navigate('/', { replace: true })
        throw e
      })
  }

  async function updateTraining(training: Training) {
    trainingApi
      .editTraining({ id: params.id, training })
      .then((res) => {
        showToast(t('training.edited'))
        setEditTraining(false)
        updateTrainintDetails(res)
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        let msg = t('server.error')
        if (e.response?.status === 404) {
          msg = t('training.not.found')
        }
        showToast(msg, ToastMode.ERROR)
        setEditTraining(false)
      })
  }

  async function deleteTraining(id: string) {
    await trainingApi
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
    <Show when={!training.loading} fallback={<Spinner remSize={8} />}>
      <Show
        when={editTraining()}
        fallback={
          <TrainingPreviewPage
            training={training()!}
            leftHeaderComponent={() => (
              <i
                class="fa-solid fa-pen fa-xl text-sky-900"
                onClick={() => setEditTraining(true)}
              ></i>
            )}
            rightHeaderComponent={() => (
              <div class="text-right">
                <ConfirmationModal
                  icon="fa-trash"
                  iconColor="text-red-500"
                  message={t('confirm.training.delete.message')}
                  confirmLabel={t('confirm.delete.training')}
                  cancelLabel={t('no.cancel')}
                  onConfirm={() => {
                    deleteTraining(params.id)
                  }}
                  onCancel={() => {}}
                />
              </div>
            )}
          />
        }
      >
        <EditTrainingPage
          saveToLocalStorage={false}
          training={createStore(JSON.parse(JSON.stringify(training())))}
          onSubmit={(tr) => updateTraining(tr as Training)}
          onBack={() => setEditTraining(false)}
        />
      </Show>
    </Show>
  )
}

export default TrainingPage
