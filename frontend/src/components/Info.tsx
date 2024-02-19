import { Component } from 'solid-js'

type MessageProps = {
  message: string
  type: 'error' | 'info'
}

const Message: Component<MessageProps> = (props) => {
  return (
    <div
      classList={{
        'bg-sky-200': props.type === 'info',
        'bg-red-300': props.type === 'error',
      }}
      class="m-4 flex items-center justify-start rounded p-4 font-bold"
    >
      {props.message}
    </div>
  )
}

export default Message
