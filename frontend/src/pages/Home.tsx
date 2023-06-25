import { Component, Show, For } from 'solid-js'
import { useTrainingsDetails } from '../state/trainings'
import DetailCard from '../components/DetailCard'
import { useNavigate } from '@solidjs/router'
import { Trans } from '@mbarzda/solid-i18next'

const Home: Component = () => {
  const [details] = useTrainingsDetails()
  const navigate = useNavigate()

  return (
    <div class="mx-auto mt-4 h-full">
      <h1 class="mx-4 text-2xl font-bold">
        <Trans key="this.weeks.trainings" />
      </h1>

      <Show
        when={!details.error}
        fallback={
          <div class="m-4 flex items-center justify-start rounded bg-red-300 p-4 font-bold">
            <Trans key="couldnt.load.trainings" />
          </div>
        }
      >
        <Show when={details()?.details?.length === 0}>
          <div class="m-4 flex items-center justify-start rounded bg-blue-200 p-4 font-bold">
            <Trans key="no.trainings" />
          </div>
        </Show>

        <For each={details()?.details}>
          {(detail) => (
            <div onClick={() => navigate('/training/' + detail.id)}>
              <DetailCard detail={detail} />
            </div>
          )}
        </For>
        {/* Add space at the bottom, so the buttons dont hide block form */}
        <div class="h-32 w-full"></div>
      </Show>
    </div>
  )
}

export default Home
