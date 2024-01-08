import { Component, Show, For } from 'solid-js'
import { useTrainingsDetailsThisWeek } from '../state/trainings'
import DetailCard from '../components/DetailCard'
import { useNavigate } from '@solidjs/router'
import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import Message from '../components/common/Info'

const Home: Component = () => {
  const [t] = useTransContext()
  const [details] = useTrainingsDetailsThisWeek()
  const navigate = useNavigate()

  return (
    <div class="h-full">
      <h1 class="mx-4 text-2xl font-bold">
        <Trans key="this.weeks.trainings" />
      </h1>

      <Show
        when={!details.error}
        fallback={
          <Message type="error" message={t('couldnt.load.trainings')} />
        }
      >
        <Show when={details()?.details?.length === 0}>
          <Message type="info" message={t('no.trainings')} />
        </Show>

        <For each={details()?.details}>
          {(detail) => (
            <div onClick={() => navigate('/training/' + detail.id)}>
              <DetailCard detail={detail} />
            </div>
          )}
        </For>
      </Show>
    </div>
  )
}

export default Home
