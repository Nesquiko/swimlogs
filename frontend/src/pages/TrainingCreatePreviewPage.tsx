import { Trans } from '@mbarzda/solid-i18next'
import { Component } from 'solid-js'
import { useShownComponent } from '../components/ShownComponentContextProvider'
import TrainingPreview from '../components/TrainingPreview'
import { useStateContext } from '../components/TrainingStateContext'

export const TrainingCreatePreviewPage: Component = () => {
  const [{ training }] = useStateContext()
  const [, setCurrentComponent] = useShownComponent()

  return (
    <div>
      <TrainingPreview training={training} />
      <button
        class="fixed bottom-0 left-4 mx-auto my-4 w-1/4 rounded border bg-purple-dark py-2 text-xl font-bold text-white"
        onClick={() => setCurrentComponent((c) => c - 1)}
      >
        <Trans key="previous" />
      </button>

      <button
        class="fixed bottom-0 right-4 mx-auto my-4 w-1/4 rounded border bg-green-600 py-2 text-xl font-bold text-white"
        onClick={() => submitTraining()}
      >
        <Trans key="finish" />
      </button>
    </div>
  )
}
