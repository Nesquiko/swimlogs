import { Component, createSignal } from 'solid-js'
import InlineDatepicker from '../components/InlineDatepicker'
import { NewTraining } from '../generated'

interface EditTrainingSessionPageProps {
  onSubmit: (training: NewTraining) => void
}

const EditTrainingSessionPage: Component<EditTrainingSessionPageProps> = () => {
  const [date, setDate] = createSignal(new Date())

  return (
    <div class="px-4">
      <InlineDatepicker onChange={setDate} />
      <pre>{JSON.stringify(date(), null, 2)}</pre>
    </div>
  )
}

export default EditTrainingSessionPage
