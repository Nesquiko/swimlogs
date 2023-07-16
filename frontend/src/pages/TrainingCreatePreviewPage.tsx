import { Trans } from '@mbarzda/solid-i18next'
import { Component } from 'solid-js'
import { useShownComponent } from '../components/ShownComponentContextProvider'
import TrainingPreview from '../components/TrainingPreview'
import { useStateContext } from '../components/TrainingStateContext'

interface TrainingCreatePreviewPageProps {
  onSubmit: () => void
}

export const TrainingCreatePreviewPage: Component<
  TrainingCreatePreviewPageProps
> = (props) => {
  const [{ training }] = useStateContext()
  const [, setCurrentComponent] = useShownComponent()

  return (
    <div>
      <TrainingPreview training={training} />
      <button
        class="fixed bottom-0 left-4 my-4 w-20 rounded-lg bg-sky-500 py-2 text-xl font-bold text-white"
        onClick={() => setCurrentComponent((c) => c - 1)}
      >
        <Trans key="previous" />
      </button>

      <button
        class="fixed bottom-0 right-4 my-4 w-28 rounded-lg bg-green-500 py-2 text-xl font-bold text-white"
        onClick={() => props.onSubmit()}
      >
        <Trans key="finish" />
      </button>
    </div>
  )
}
