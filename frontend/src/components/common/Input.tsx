import { Component, createSignal, onMount, Show } from 'solid-js'
import { randomId } from '../../lib/str'

type InputProps<T> = {
  onChange: (value: T) => void
  value?: T
  label?: string
  placeholder?: string
  maxLength?: number
  error?: string
}

const MIN_TEXTAREA_HEIGHT = 100

const TextAreaInput: Component<InputProps<string>> = (props) => {
  const id = randomId()
  let textAreaRef: HTMLTextAreaElement
  const [textAreaHeight, setTextAreaHeight] = createSignal(MIN_TEXTAREA_HEIGHT)

  onMount(() => {
    changeHeight(textAreaRef.scrollHeight)
  })

  const changeHeight = (newHeight: number) => {
    if (newHeight < MIN_TEXTAREA_HEIGHT) {
      newHeight = MIN_TEXTAREA_HEIGHT
    }
    setTextAreaHeight(newHeight)
  }

  return (
    <div>
      <Show when={props.label}>
        <label
          classList={{
            'text-red-500 ': props.error !== undefined,
          }}
          class="text-xl"
          for={id}
        >
          {props.label}
        </label>
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
          changeHeight(textAreaRef.scrollHeight)
        }}
      />
      <ErrorMessage error={props.error} />
    </div>
  )
}

const NumberInput: Component<InputProps<number>> = (props) => {
  const id = randomId()

  return (
    <div class="">
      <Show when={props.label}>
        <label
          classList={{
            'text-red-500 ': props.error !== undefined,
          }}
          class="text-xl md:pr-8"
          for={id}
        >
          {props.label}
        </label>
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
        class="w-full rounded-lg border p-2 text-center text-lg focus:outline-none md:w-min"
        onChange={(e) => props.onChange(parseInt(e.target.value))}
      />

      <ErrorMessage error={props.error} />
    </div>
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

export { TextAreaInput, NumberInput }
