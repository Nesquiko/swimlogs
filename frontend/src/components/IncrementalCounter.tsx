import { Component, Show } from 'solid-js';
import { randomId } from '../lib/str';

interface CounterProps {
  onChange: (n: number) => void;
  value: number;

  min?: number;
  max?: number;
  step?: number;
  label?: string;
  error?: boolean;
}

const IncrementalCounter: Component<CounterProps> = (props) => {
  const id = randomId();
  const step = props.step ?? 1;

  const isLessEqualThanMin = () =>
    props.min !== undefined ? props.value <= props.min : false;
  const isMoreEqualThanMax = () =>
    props.max !== undefined ? props.value >= props.max : false;

  return (
    <div class="flex items-center justify-between md:block md:w-44">
      <Show when={props.label}>
        <label
          for={id}
          classList={{ 'text-red-500': props.error }}
          class="block text-center text-xl"
        >
          {props.label}
        </label>
      </Show>
      <div
        classList={{ 'border-red-500': props.error }}
        class="flex h-10 rounded-lg"
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
          class="w-14 cursor-pointer rounded-l-lg border-r-0 border-slate-300 leading-10"
          onClick={() => props.onChange(props.value - step)}
        >
          <i class="fa-solid fa-minus"></i>
        </button>
        <input
          id={id}
          type="number"
          inputmode="numeric"
          classList={{ 'text-red-500': props.error }}
          class="w-16 p-2 text-center text-lg border-x-0 border-y border-slate-200 focus:outline focus:outline-sky-300"
          min={props.min}
          max={props.max}
          value={props.value}
          onChange={(e) => {
            let n = parseInt(e.target.value);
            const isLessThanMin =
              props.min !== undefined ? n < props.min : false;
            const isMoreThanMax =
              props.max !== undefined ? n > props.max : false;

            if (Number.isNaN(n)) {
              n = props.value ?? 0;
            } else if (isLessThanMin) {
              n = props.min ?? 0;
            } else if (isMoreThanMax) {
              n = props.max ?? 0;
            }

            props.onChange(n);
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
          class="w-14 cursor-pointer rounded-r-lg border-l-0 border-slate-300 leading-10"
          onClick={() => props.onChange(props.value + step)}
        >
          <i class="fa-solid fa-plus"></i>
        </button>
      </div>
    </div>
  );
};

export default IncrementalCounter;
