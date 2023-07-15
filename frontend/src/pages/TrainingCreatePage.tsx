import { useNavigate } from '@solidjs/router'
import { Component, createEffect, createResource, createSignal } from 'solid-js'
import { createStore } from 'solid-js/store'
import { Dynamic } from 'solid-js/web'
import { NewTraining, ResponseError, Session } from '../generated'
import { addTrainingDetail, trainingApi } from '../state/trainings'
import { NullDateTime, PAGE_SIZE } from '../lib/consts'
import { useTransContext } from '@mbarzda/solid-i18next'
import { openToast, ToastType } from '../components/Toast'
import TrainingSetsForm from './TrainingSetsForm'
import { TrainingSessionForm } from './TrainingSessionForm'
import { TrainingCreatePreviewPage } from './TrainingCreatePreviewPage'
import { TrainingStateContextProvider } from '../components/TrainingStateContext'
import { ShownComponentContextProvider } from '../components/ShownComponentContextProvider'
import { SessionsContextProvider } from '../components/SessionsContextProvider'
import { sessionApi } from '../state/session'

const TrainingCreatePage: Component = () => {
  const [training, setTraining] = createStore<NewTraining>({
    start: NullDateTime,
    durationMin: 60,
    totalDistance: 100,
    sets: []
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
        console.error('error', e)
        openToast(
          t('training.creation.error', 'Error creating training'),
          ToastType.ERROR
        )
        navigate('/', { replace: true })
      })
  }

  const [sessionsPage, setSessionsPage] = createSignal(0)
  const [totalSessions, setTotalSessions] = createSignal(0)
  const isLastPage = () => (sessionsPage() + 1) * PAGE_SIZE >= totalSessions()
  const [sessions] = createResource(sessionsPage, getSessions)
  const [serverError, setServerError] = createSignal<string | undefined>(
    undefined
  )

  async function getSessions(page: number): Promise<Session[]> {
    return sessionApi
      .getSessions({ page, pageSize: PAGE_SIZE })
      .then((res) => {
        setTotalSessions(res.pagination.total)
        return res.sessions
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        setServerError(e.message)
        return Promise.resolve([])
      })
  }

  const [currentComponent, setCurrentComponent] = createSignal(0)
  const comps = [
    TrainingSessionForm,
    TrainingSetsForm,
    TrainingCreatePreviewPage
  ]

  return (
    <SessionsContextProvider
      sessions={sessions}
      fetchNextSessionPage={() => setSessionsPage(sessionsPage() + 1)}
      fetchPrevSessionPage={() => setSessionsPage(sessionsPage() - 1)}
      page={sessionsPage}
      isLastPage={isLastPage}
      serverError={serverError}
    >
      <ShownComponentContextProvider
        currentComponentSignal={[currentComponent, setCurrentComponent]}
      >
        <TrainingStateContextProvider newTraining={training}>
          <Dynamic
            {...{ onSubmit: createTraining }}
            component={comps[currentComponent()]}
          />
        </TrainingStateContextProvider>
      </ShownComponentContextProvider>
    </SessionsContextProvider>
  )
}

export default TrainingCreatePage
