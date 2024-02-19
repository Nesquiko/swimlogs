import { Trans } from '@mbarzda/solid-i18next'
import { useNavigate } from '@solidjs/router'
import { Component, createResource, createSignal, For, Show } from 'solid-js'
import Pagination from '../components/Pagination'
import { ResponseError, TrainingDetail } from 'swimlogs-api'
import { trainingApi } from '../state/trainings'
import { formatDate } from '../lib/datetime'

const PAGE_SIZE = 8

const TrainingHistoryPage: Component = () => {
  const navigate = useNavigate()

  const [detailsPage, setDetailsPage] = createSignal(0)
  const [totalDetails, setTotalDetails] = createSignal(0)
  const [serverError, setServerError] = createSignal(false)
  const isLastPage = () => (detailsPage() + 1) * PAGE_SIZE >= totalDetails()

  const cachedTrainingDetails = new Map<number, TrainingDetail[]>()
  const [details] = createResource(detailsPage, getTrainingDetals)
  async function getTrainingDetals(page: number): Promise<TrainingDetail[]> {
    if (cachedTrainingDetails.has(page)) {
      setServerError(false)
      return Promise.resolve(cachedTrainingDetails.get(page)!)
    }

    return trainingApi
      .trainingDetails({ page, pageSize: PAGE_SIZE })
      .then((res) => {
        setTotalDetails(res.pagination.total)
        setServerError(false)
        cachedTrainingDetails.set(page, res.details)
        return res.details
      })
      .catch((e: ResponseError) => {
        console.error('error', e)
        setServerError(true)
        return Promise.resolve([])
      })
  }

  return (
    <div class="h-screen">
      <h1 class="m-4 text-2xl font-bold">
        <Trans key="trainings.history" />
      </h1>
      <div class="h-5/6">
        <Show when={serverError()}>
          <div class="m-4 flex items-center justify-start rounded bg-red-300 p-4 font-bold">
            <Trans key="couldnt.load.trainings" />
          </div>
        </Show>
        <For each={details()}>
          {(detail) => (
            <div onClick={() => navigate('/training/' + detail.id)}>
              <DetailCard detail={detail} />
            </div>
          )}
        </For>
      </div>
      <Pagination
        prevDisabled={detailsPage() === 0}
        onPrevPage={() => setDetailsPage((i) => i - 1)}
        nextDisabled={isLastPage() || serverError()}
        onNextPage={() => setDetailsPage((i) => i + 1)}
      />
    </div>
  )
}

interface DetailProps {
  detail: TrainingDetail
}

const DetailCard: Component<DetailProps> = (props) => {
  const day = props.detail.start
    .toLocaleString('en', { weekday: 'long' })
    .toLowerCase()
  return (
    <div class="z-100 mx-auto my-4 w-11/12 cursor-pointer rounded-lg border border-solid border-slate-200 bg-white p-2 shadow">
      <h2 class="flex justify-between text-left text-xl">
        <b class="w-1/3">
          <Trans key={day} />
        </b>
        <p>{formatDate(props.detail.start)}</p>
      </h2>
      <div class="flex justify-between">
        <p class="text-base">{props.detail.durationMin} min</p>
        <p class="text-base">{props.detail.totalDistance}m</p>
      </div>
    </div>
  )
}

export default TrainingHistoryPage
