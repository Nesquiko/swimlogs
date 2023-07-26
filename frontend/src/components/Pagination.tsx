import { Component } from 'solid-js'

interface PaginationProps {
  onPrevPage: () => void
  onNextPage: () => void
  prevDisabled: boolean
  nextDisabled: boolean
}

const Pagination: Component<PaginationProps> = (props) => {
  return (
    <div class="space-x-8 p-2 text-center">
      <button
        classList={{ 'text-slate-300 pointer-events-none': props.prevDisabled }}
        disabled={props.prevDisabled}
        class="inline-flex cursor-pointer text-black"
        onClick={() => props.onPrevPage()}
      >
        <svg
          aria-hidden="true"
          class="h-8 w-8"
          fill="currentColor"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fill-rule="evenodd"
            d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
            clip-rule="evenodd"
          ></path>
        </svg>
      </button>
      <button
        classList={{ 'text-slate-300 pointer-events-none': props.nextDisabled }}
        class="inline-flex cursor-pointer text-black"
        disabled={props.nextDisabled}
        onClick={() => props.onNextPage()}
      >
        <svg
          aria-hidden="true"
          class="h-8 w-8"
          fill="currentColor"
          viewBox="0 0 20 20"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path
            fill-rule="evenodd"
            d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
            clip-rule="evenodd"
          ></path>
        </svg>
      </button>
    </div>
  )
}

export default Pagination
