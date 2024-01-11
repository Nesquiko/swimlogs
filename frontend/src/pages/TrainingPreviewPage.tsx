import { Trans } from '@mbarzda/solid-i18next'
import { useNavigate } from '@solidjs/router'
import { Component } from 'solid-js'
import { createStore } from 'solid-js/store'
import { NewTraining } from '../generated'

const TrainingPreviewPage: Component = () => {
  const navigate = useNavigate()
  const [training, setTraining] = createStore<NewTraining>({
    start: new Date(),
    durationMin: 0,
    totalDistance: 0,
    sets: [],
  })

  return (
    <div class="px-4">
      <h1 class="flex justify-between text-2xl font-bold">
        <Trans key="total.distance" />
        <span>{training.totalDistance / 100}km</span>
      </h1>

      <button
        class="fixed bottom-2 right-2 h-16 w-16 rounded-lg bg-sky-500"
        onClick={() => navigate('/training/set/new')}
      >
        <i class="fa-solid fa-plus fa-2xl text-white"></i>
      </button>
    </div>
  )
}

export default TrainingPreviewPage
