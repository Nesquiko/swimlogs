import { Component } from 'solid-js'
import { TrainingDetail } from '../generated'

interface DetailProps {
  detail: TrainingDetail
}

const DetailCard: Component<DetailProps> = (props) => {
  return (
    <div class="mx-auto my-2 rounded-lg border border-solid border-slate-100 bg-white p-2 shadow">
      <div class="flex justify-between">
        <p class="text-base">
          <b>Start:</b> {props.detail.startTime}
        </p>
        <p class="text-base">
          <b>Duration:</b> {props.detail.durationMin} min
        </p>
        <p class="text-base">
          <b>Distance:</b> {props.detail.totalDist}m
        </p>
      </div>
    </div>
  )
}

export default DetailCard
