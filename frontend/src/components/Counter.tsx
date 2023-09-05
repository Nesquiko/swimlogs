import { Component, createEffect, createSignal } from 'solid-js'

interface CounterProps {
  onChange: (n: number) => void
  initial: number
  min?: number
  max?: number
}

const Counter: Component<CounterProps> = (props) => {
  const [number, setNumber] = createSignal(props.initial)

  const changeNumber = (n: number) => {
    setNumber(n)
    props.onChange(number())
  }

  createEffect(() => setNumber(props.initial)) // updates number when on repeated open

  const isLessEqualThanMin = () => (props.min ? number() <= props.min! : false)
  const isLessThanMin = () => (props.min ? number() < props.min! : false)
  const isMoreEqualThanMax = () => (props.max ? number() >= props.max! : false)
  const isMoreThanMax = () => (props.max ? number() > props.max! : false)

  return (
    <div class="relative flex h-10 max-w-fit flex-row rounded-lg border border-slate-300">
      <button
        disabled={isLessEqualThanMin()}
        classList={{
          'text-slate-300 pointer-events-none bg-sky-100': isLessEqualThanMin(),
          'bg-sky-400': !isLessEqualThanMin()
        }}
        class="h-full w-8 cursor-pointer rounded-l border-r border-slate-300"
        onClick={() => changeNumber(number() - 1)}
      >
        <span class="m-auto text-2xl">-</span>
      </button>
      <input
        type="number"
        placeholder="1"
        class="w-20  p-2 text-center text-lg focus:outline-none"
        value={number()}
        onChange={(e) => {
          let n = parseInt(e.target.value)
          if (Number.isNaN(n) || isLessThanMin() || isMoreThanMax()) {
            n = props.min ?? props.initial ?? 0
          }
          changeNumber(n)
        }}
      />
      <button
        disabled={isMoreEqualThanMax()}
        classList={{
          'text-slate-300 pointer-events-none bg-sky-100': isMoreEqualThanMax(),
          'bg-sky-400': !isMoreEqualThanMax()
        }}
        class="h-full w-8 cursor-pointer rounded-r border-l bg-sky-400"
        onClick={() => changeNumber(number() + 1)}
      >
        <span class="m-auto border-slate-300 text-2xl">+</span>
      </button>
    </div>
  )
}

export default Counter
