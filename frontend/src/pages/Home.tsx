import { Component, Show, For } from 'solid-js'
import { useTrainingsDetailsThisWeek } from '../state/trainings'
import { useNavigate } from '@solidjs/router'
import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import Message from '../components/common/Info'
import { Day, TrainingDetail } from '../generated'
import { dateForDayOfWeek, formatDate, formatTime } from '../lib/datetime'

const Home: Component = () => {
  const [t] = useTransContext()
  const [details] = useTrainingsDetailsThisWeek()
  const navigate = useNavigate()

  const detailItem = (detail: TrainingDetail) => {
    return (
      <div
        onClick={() => navigate('/training/' + detail.id)}
        class="flex cursor-pointer justify-between rounded-lg border border-solid border-slate-200 px-8 py-2"
      >
        <span class="text-lg">{formatTime(detail.start)}</span>
        <span class="text-lg">{detail.totalDistance / 1000}km</span>
      </div>
    )
  }

  const sortedTrainingsPerDay = (details: TrainingDetail[]) => {
    if (details.length === 0) return <></>

    const dayNames: Day[] = Object.values(Day)
    const dayToDetails: {
      [K in Day]: { details: TrainingDetail[]; date: Date }
    } = {
      Monday: { details: [], date: dateForDayOfWeek(1) },
      Tuesday: { details: [], date: dateForDayOfWeek(2) },
      Wednesday: { details: [], date: dateForDayOfWeek(3) },
      Thursday: { details: [], date: dateForDayOfWeek(4) },
      Friday: { details: [], date: dateForDayOfWeek(5) },
      Saturday: { details: [], date: dateForDayOfWeek(6) },
      Sunday: { details: [], date: dateForDayOfWeek(7) },
    }

    for (const detail of details) {
      const day = dayNames[(detail.start.getDay() + 6) % 7]
      dayToDetails[day].details.push(detail)
    }

    return (
      <For each={Object.entries(dayToDetails)}>
        {(entry) => {
          const day = entry[0].toLowerCase()
          return (
            <div class="py-4">
              <div class="flex justify-between text-xl font-bold">
                <span>{t(day)}</span>
                <span>{formatDate(entry[1].date)}</span>
              </div>
              <Show
                when={entry[1].details.length !== 0}
                fallback={
                  <div class="text-center">
                    <i class="fa-solid fa-minus"></i>{' '}
                    <i class="fa-solid fa-minus"></i>
                  </div>
                }
              >
                <div class="space-y-2">
                  <For each={entry[1].details}>{detailItem}</For>
                </div>
              </Show>
            </div>
          )
        }}
      </For>
    )
  }

  return (
    <div class="h-full px-4 pb-28">
      <button
        class="fixed bottom-2 right-2 h-16 w-16 rounded-lg bg-sky-500"
        onClick={() => navigate('/training/new')}
      >
        <i class="fa-solid fa-plus fa-2xl text-white"></i>
      </button>

      <Show
        when={!details.error}
        fallback={
          <Message type="error" message={t('couldnt.load.trainings')} />
        }
      >
        <h1 class="text-2xl font-bold">
          <Trans key="this.week" />
        </h1>
        <div class="flex justify-between">
          <Trans key="total.distance.swam" />
          <p>{(details()?.distance ?? 0) / 1000}km</p>
        </div>
        <Show when={details()?.details?.length === 0}>
          <Message type="info" message={t('no.trainings')} />
        </Show>
        {sortedTrainingsPerDay(details()?.details ?? [])}
      </Show>
    </div>
  )
}

export default Home
