import { Component, createSignal, For, Show } from 'solid-js'
import { randomId } from '../../lib/str'

type InputProps<T> = {
  onChange: (value: T) => void
  value?: T
  label?: string
  placeholder?: string
  maxLength?: number
  error?: string
}

const INITIAL_HEIGHT = 44
const LINE_HEIGHT = 28 // tailwind 'text-lg' line height is 28px

const TextAreaInput: Component<InputProps<string>> = (props) => {
  const id = randomId()
  let textAreaRef: HTMLTextAreaElement
  const [textAreaHeight, setTextAreaHeight] = createSignal(INITIAL_HEIGHT)

  const changeHeight = () => {
    const newLines = props.value?.match(/\n/g)?.length ?? 0
    const newHeightWithLines = INITIAL_HEIGHT + newLines * LINE_HEIGHT

    setTextAreaHeight(newHeightWithLines)
  }
  changeHeight()

  return (
    <div>
      <Show when={props.label}>
        <Label id={id} label={props.label!} error={props.error} />
      </Show>
      <textarea
        ref={textAreaRef!}
        id={id}
        style={{ height: textAreaHeight() + 'px' }}
        classList={{
          'border-red-500 text-red-500 focus:border-red-500':
            props.error !== undefined,
          'border-slate-300 focus:border-sky-500': props.error === undefined,
        }}
        class="w-full overflow-y-hidden rounded-lg border p-2 text-lg focus:outline-none"
        value={props.value ?? ''}
        placeholder={props.placeholder ?? ''}
        maxlength={props.maxLength}
        onInput={(e) => {
          props.onChange(e.currentTarget.value)
          changeHeight()
        }}
      />
      <ErrorMessage error={props.error} />
    </div>
  )
}

const NumberInput: Component<InputProps<number>> = (props) => {
  const id = randomId()

  return (
    <div class="md:w-96">
      <Show when={props.label}>
        <Label id={id} label={props.label!} error={props.error} />
      </Show>
      <input
        id={id}
        type="number"
        value={props.value ?? ''}
        placeholder={props.placeholder ?? ''}
        classList={{
          'border-red-500 text-red-500 focus:border-red-500':
            props.error !== undefined,
          'border-slate-300 focus:border-sky-500': props.error === undefined,
        }}
        class="w-full rounded-lg border p-2 text-center text-lg focus:outline-none md:float-right md:w-44"
        onChange={(e) => props.onChange(parseInt(e.target.value))}
      />

      <ErrorMessage error={props.error} />
    </div>
  )
}

type SelectOption<T> = {
  label: string
  value: T
}

type SelectInputProps<T> = {
  onChange: (value: SelectOption<T> | undefined) => void
  options: SelectOption<T>[]
  value?: T
  label?: string
  noneOption?: string
  error?: string
}

const SelectInput = <T extends object>(props: SelectInputProps<T>) => {
  const id = randomId()
  const [selectedOption, setSelectedOption] = createSignal<
    SelectOption<T> | undefined
  >(undefined)

  return (
    <div class="w-44 bg-yellow-200 py-2">
      <div class="block">
        <Show when={props.label}>
          <Label id={id} label={props.label!} error={props.error} />
        </Show>
      </div>
      <select
        id={id}
        class="w-44 rounded-lg border border-solid border-slate-300 bg-white p-2 text-start text-lg focus:outline-none"
        classList={{
          'border-red-500 text-red-500 focus:border-red-500':
            props.error !== undefined,
          'border-slate-300 focus:border-sky-500': props.error === undefined,
        }}
        onChange={(e) => {
          let option: SelectOption<T> | undefined
          if (e.target.value === undefined) {
            option = undefined
          } else {
            option = props.options.find((o) => o.label === e.target.value)
          }
          setSelectedOption(option)
          props.onChange(option)
        }}
      >
        <Show when={props.noneOption}>
          <option selected={selectedOption() === undefined} value={undefined}>
            {props.noneOption}
          </option>
        </Show>
        <For each={props.options}>
          {(opt) => (
            <option
              selected={opt.label === selectedOption()?.label}
              value={opt.label}
            >
              {opt.label}
            </option>
          )}
        </For>
      </select>
      <ErrorMessage error={props.error} />
    </div>
  )
}

type LabelProps = {
  id: string
  label: string
  error?: string
}

const Label: Component<LabelProps> = (props) => {
  return (
    <label
      classList={{
        'text-red-500 ': props.error !== undefined,
      }}
      class="text-xl md:pr-8"
      for={props.id}
    >
      {props.label}
    </label>
  )
}

type ErrorMessageProps = {
  error?: string
}

const ErrorMessage: Component<ErrorMessageProps> = (props) => {
  return (
    <div class="h-6 text-red-500">
      <Show when={props.error}>
        <i class="fa-solid fa-circle-exclamation"></i>
        <span class="px-2">{props.error}</span>
      </Show>
    </div>
  )
}

export { TextAreaInput, NumberInput, SelectInput }
