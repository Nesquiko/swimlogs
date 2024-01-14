import { Component, createSignal, For, Show } from 'solid-js'
import { randomId } from '../../lib/str'

type TextAreaProps = {
  onInput: (value: string) => void
  value?: string
  label?: string
  placeholder?: string
  maxLength?: number

  validated?: boolean
  error?: string
}

const LINE_HEIGHT = 24 // tailwind 'text-base' line height is 24px
const INITIAL_HEIGHT = 2 * LINE_HEIGHT

const TextAreaInput: Component<TextAreaProps> = (props) => {
  const id = randomId()
  let textAreaRef: HTMLTextAreaElement

  const calculateHeight = () => {
    const newLines = props.value?.match(/\n/g)?.length ?? 0
    const newHeightWithLines = INITIAL_HEIGHT + newLines * LINE_HEIGHT
    return newHeightWithLines
  }
  const [textAreaHeight, setTextAreaHeight] = createSignal(
    props.value ? calculateHeight() : INITIAL_HEIGHT
  )

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
        class="w-full overflow-y-hidden rounded-lg border p-2 text-base focus:outline-none"
        value={props.value ?? ''}
        placeholder={props.placeholder ?? ''}
        maxlength={props.maxLength}
        onInput={(e) => {
          props.onInput(e.currentTarget.value)
          setTextAreaHeight(calculateHeight())
        }}
      />
      <ErrorMessage error={props.error} />
    </div>
  )
}

type NumberInputProps = {
  onChange: (n: number | undefined) => void
  value?: number
  label?: string
  placeholder?: string
  error?: string
}

const NumberInput: Component<NumberInputProps> = (props) => {
  const id = randomId()

  return (
    <div class="w-full md:w-96">
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
        class="w-full rounded-lg border p-2 text-center text-lg shadow focus:outline-none md:float-right md:w-44"
        onChange={(e) => {
          const n = parseInt(e.target.value)
          if (isNaN(n)) {
            props.onChange(undefined)
            return
          }
          props.onChange(parseInt(e.target.value))
        }}
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
  validated?: boolean
  error?: string
}

const SelectInput = <T extends object>(props: SelectInputProps<T>) => {
  const id = randomId()
  const [selectedOption, setSelectedOption] = createSignal<
    SelectOption<T> | undefined
  >(undefined)

  return (
    <div class="w-44 py-2">
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
      <Show when={!!props.validated}>
        <ErrorMessage error={props.error} />
      </Show>
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
