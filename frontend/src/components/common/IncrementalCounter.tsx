import { Component, Show } from 'solid-js'
import { randomId } from '../../lib/str'

interface CounterProps {
  onChange: (n: number) => void
  value: number

  min?: number
  max?: number
  label?: string
  error?: boolean
}

const IncrementalCounter: Component<CounterProps> = (props) => {
  const id = randomId()

  const isLessEqualThanMin = () =>
    props.min !== undefined ? props.value <= props.min : false
  const isMoreEqualThanMax = () =>
    props.max !== undefined ? props.value >= props.max : false

  return (
    <div class="flex items-center justify-between md:w-96">
      <Show when={props.label}>
        <label
          for={id}
          classList={{ 'text-red-500': props.error }}
          class="text-center text-xl"
        >
          {props.label}
        </label>
      </Show>
      <div
        classList={{ 'border-red-500': props.error }}
        class="flex h-10 flex-row rounded-lg border border-slate-300"
      >
        <button
          disabled={isLessEqualThanMin()}
          classList={{
            'bg-sky-400': !props.error && !isLessEqualThanMin(),
            'text-slate-300 pointer-events-none select-none bg-sky-100':
              !props.error && isLessEqualThanMin(),
            'bg-red-400': props.error && !isLessEqualThanMin(),
            'text-slate-300 pointer-events-none select-none bg-red-100':
              props.error && isLessEqualThanMin(),
          }}
          class="h-full w-12 cursor-pointer rounded-l-lg border-r border-slate-300"
          onClick={() => props.onChange(props.value - 1)}
        >
          <span class="m-auto text-2xl">-</span>
        </button>
        <input
          id={id}
          type="number"
          classList={{ 'text-red-500': props.error }}
          class="w-14 p-2 text-center text-lg [appearance:textfield] focus:outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none"
          min={props.min}
          max={props.max}
          value={props.value}
          onChange={(e) => {
            let n = parseInt(e.target.value)
            const isLessThanMin =
              props.min !== undefined ? n < props.min : false
            const isMoreThanMax =
              props.max !== undefined ? n > props.max : false

            if (Number.isNaN(n)) {
              n = props.value ?? 0
            } else if (isLessThanMin) {
              n = props.min ?? 0
            } else if (isMoreThanMax) {
              n = props.max ?? 0
            }

            props.onChange(n)
          }}
        />
        <button
          disabled={isMoreEqualThanMax()}
          classList={{
            'bg-sky-400': !props.error && !isMoreEqualThanMax(),
            'text-slate-300 pointer-events-none select-none bg-sky-100':
              !props.error && isMoreEqualThanMax(),
            'bg-red-400': props.error && !isMoreEqualThanMax(),
            'text-slate-300 pointer-events-none select-none bg-red-100':
              props.error && isMoreEqualThanMax(),
          }}
          class="h-full w-12 cursor-pointer rounded-r-lg border-l bg-sky-400"
          onClick={() => props.onChange(props.value + 1)}
        >
          <span class="m-auto border-slate-300 text-2xl">+</span>
        </button>
      </div>
    </div>
  )
}

export default IncrementalCounter
