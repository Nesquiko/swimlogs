import { Component } from 'solid-js'
import { createStore } from 'solid-js/store'
import { Button } from '@suid/material'
import { InvalidTraining, Training } from '../generated'
import TraningForm from '../component/TrainingFrom'

const [training, setTraining] = createStore<Training>({
  id: 'id',
  date: new Date(),
  blocks: [],
  version: 0
})

const [invalidTraining, setInvalidTraining] = createStore<InvalidTraining>({})

const validateTraining = (t: Training) => {
  !t.day
    ? setInvalidTraining('day', "Day can't be empty")
    : setInvalidTraining('day', undefined)

  if (t.durationMin === undefined)
    setInvalidTraining('durationMin', "Duration can't be empty")
  else if (t.durationMin < 0)
    setInvalidTraining('durationMin', 'Duration must be positive')
  else setInvalidTraining('durationMin', undefined)
}

const CreateTraining: Component = () => {
  return (
    <div class="mx-auto w-11/12 text-center">
      <h1 class="my-4 text-2xl font-bold">Create Training</h1>

      <TraningForm
        training={training}
        setTraining={setTraining}
        invalidTraining={invalidTraining}
      />
      <pre>Training: {JSON.stringify(training, null, 4)}</pre>

      <Button
        onClick={() => validateTraining(training)}
        variant="contained"
        size="large"
      >
        Create
      </Button>
    </div>
  )
}

export default CreateTraining
