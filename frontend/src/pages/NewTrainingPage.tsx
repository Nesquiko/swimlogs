import { useTransContext } from '@mbarzda/solid-i18next'
import { useNavigate } from '@solidjs/router'
import { Component } from 'solid-js'
import { createStore } from 'solid-js/store'
import { showToast } from '../App'
import { ToastMode } from '../components/common/DismissibleToast'
import { NewTraining } from '../generated'
import { loadTrainingFromLocalStorage } from '../state/local-storage'
import { addTrainingDetail, trainingApi } from '../state/trainings'
import EditTrainingPage from './EditTrainingPage'

const NewTrainingPage: Component = () => {
  const [t] = useTransContext()
  const navigate = useNavigate()
  const defaultStart = new Date()
  defaultStart.setHours(18, 0, 0, 0)
  const training = createStore<NewTraining>(
    loadTrainingFromLocalStorage() ?? {
      start: defaultStart,
      durationMin: 60,
      totalDistance: 0,
      sets: [],
    }
  )

  const onSubmit = (training: NewTraining) => {
    trainingApi
      .createTraining({ newTraining: training })
      .then((res) => {
        addTrainingDetail(res)
        showToast(t('training.created', 'Training created'))
      })
      .catch((e) => {
        console.error('error', e)
        showToast(t('training.creation.error'), ToastMode.ERROR)
      })
      .finally(() => {
        navigate('/', { replace: true })
      })
  }

  return (
    <EditTrainingPage
      saveToLocalStorage={true}
      training={training}
      onSubmit={onSubmit}
      onDelete={() => navigate('/', { replace: true })}
      onBack={() => navigate('/', { replace: true })}
    />
  )
}

export default NewTrainingPage
