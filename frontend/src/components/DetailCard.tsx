import { Trans } from '@mbarzda/solid-i18next'
import { Component } from 'solid-js'
import { TrainingDetail } from '../generated'
import { formatDate } from '../lib/datetime'

interface DetailProps {
  detail: TrainingDetail
}

const DetailCard: Component<DetailProps> = (props) => {
  const day = props.detail.date
    .toLocaleString('en', { weekday: 'long' })
    .toLowerCase()
  return (
    <div class="z-100 mx-auto my-4 w-11/12 cursor-pointer rounded-lg border border-solid border-slate-200 bg-white p-2 shadow">
      <h2 class="flex justify-between text-left text-xl">
        <b class="w-1/3">
          <Trans key={day} />
        </b>
        <p>{formatDate(props.detail.date)}</p>
      </h2>
      <div class="flex justify-between">
        <p>{props.detail.startTime}</p>
        <p class="text-base">{props.detail.durationMin} min</p>
        <p class="text-base">{props.detail.totalDistance}m</p>
      </div>
    </div>
  )
}

export default DetailCard
