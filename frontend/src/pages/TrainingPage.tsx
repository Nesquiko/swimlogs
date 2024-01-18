import { useTransContext } from '@mbarzda/solid-i18next'
import { useNavigate, useParams } from '@solidjs/router'
import { Component, createResource, createSignal, Show } from 'solid-js'
import { createStore } from 'solid-js/store'
import { showToast } from '../App'
import { ToastMode } from '../components/common/DismissibleToast'
import ConfirmationModal from '../components/ConfirmationModal'
import { setOnBackOverride } from '../components/Header'
import Spinner from '../components/Spinner'
import { ResponseError, Training } from '../generated'
import { removeFromTrainingsDetails, trainingApi } from '../state/trainings'
import EditTrainingPage from './EditTrainingPage'
import TrainingPreviewPage from './TrainingPreviewPage'

const TrainingPage: Component = () => {
  const params = useParams()
  const [t] = useTransContext()
  const [training, { mutate, refetch }] = createResource(
    () => params.id,
    getTraining
  )
  const [editTraining, setEditTraining] = createSignal(false)

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
                  onConfirm={() => deleteTraining(params.id)}
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
          onSubmit={(t) => {
			// TODO backend for this
            mutate(t as Training)
            setEditTraining(false)
          }}
          onBack={() => setEditTraining(false)}
        />
      </Show>
    </Show>
  )
}

export default TrainingPage
