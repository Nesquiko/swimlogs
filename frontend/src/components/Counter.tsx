import { Component, createEffect, createSignal, Show } from 'solid-js'

interface CounterProps {
  onChange: (n: number) => void
  initial: number
  min?: number
  max?: number
  label?: string
  error?: boolean
}

const Counter: Component<CounterProps> = (props) => {
  const [number, setNumber] = createSignal(props.initial)

  const changeNumber = (n: number) => {
    setNumber(n)
    props.onChange(number())
  }

  createEffect(() => setNumber(props.initial)) // updates number when on repeated open

  const isLessEqualThanMin = () =>
    props.min !== undefined ? number() <= props.min! : false
  const isMoreEqualThanMax = () =>
    props.max !== undefined ? number() >= props.max! : false

  return (
    <div>
      <Show when={props.label}>
        <p
          classList={{ 'text-red-500': props.error }}
          class="text-center text-base"
        >
          {props.label}
        </p>
      </Show>
      <div
        classList={{ 'border-red-400': props.error }}
        class="relative flex h-10 max-w-fit flex-row rounded-lg border border-slate-300"
      >
        <button
          disabled={isLessEqualThanMin()}
          classList={{
            'text-slate-300 pointer-events-none bg-sky-100':
              isLessEqualThanMin(),
            'bg-sky-400': !isLessEqualThanMin(),
          }}
          class="h-full w-8 cursor-pointer rounded-l-lg border-r border-slate-300"
          onClick={() => changeNumber(number() - 1)}
        >
          <span class="m-auto text-2xl">-</span>
        </button>
        <input
          type="number"
          placeholder="1"
          classList={{ 'text-red-500': props.error }}
          class="w-20 p-2 text-center text-lg focus:outline-none"
          value={number()}
          onChange={(e) => {
            let n = parseInt(e.target.value)
            const isLessThanMin =
              props.min !== undefined ? n < props.min! : false
            const isMoreThanMax =
              props.max !== undefined ? n > props.max! : false

            if (Number.isNaN(n)) {
              n = props.initial ?? 0
            } else if (isLessThanMin) {
              n = props.min ?? 0
            } else if (isMoreThanMax) {
              n = props.max ?? 0
            }

            changeNumber(n)
            e.target.value = n.toString()
          }}
        />
        <button
          disabled={isMoreEqualThanMax()}
          classList={{
            'text-slate-300 pointer-events-none bg-sky-100':
              isMoreEqualThanMax(),
            'bg-sky-400': !isMoreEqualThanMax(),
          }}
          class="h-full w-8 cursor-pointer rounded-r-lg border-l bg-sky-400"
          onClick={() => changeNumber(number() + 1)}
        >
          <span class="m-auto border-slate-300 text-2xl">+</span>
        </button>
      </div>
    </div>
  )
}

export default Counter
