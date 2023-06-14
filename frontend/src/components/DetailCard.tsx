import { Component } from 'solid-js'
import { TrainingDetail } from '../generated'

interface DetailProps {
  detail: TrainingDetail
}

const DetailCard: Component<DetailProps> = (props) => {
  return (
    <div class="z-100 mx-auto my-4 w-11/12 rounded-lg border border-solid border-slate-200 bg-white p-2 shadow">
      <h2 class="text-left text-xl">
        <b>{props.detail.date.toLocaleString('en', { weekday: 'long' })}</b>{' '}
        {props.detail.date.toLocaleDateString('sk-SK').replaceAll(' ', '')}
      </h2>
      <div class="flex justify-between">
        <p class="text-base">
          <b>Start:</b> {props.detail.startTime}
        </p>
        <p class="text-base">
          <b>Duration:</b> {props.detail.durationMin} min
        </p>
        <p class="text-base">
          <b>Distance:</b> {props.detail.totalDistance}m
        </p>
      </div>
    </div>
  )
}

export default DetailCard
