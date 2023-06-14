import { Component, createSignal, Show, For } from 'solid-js'
import { useTrainingsDetails } from '../state/trainings'
import DetailCard from '../components/DetailCard'
import { A, useNavigate } from '@solidjs/router'

const Home: Component = () => {
  const [details] = useTrainingsDetails()
  const [open, setOpen] = createSignal(false)
  const navigate = useNavigate()

  return (
    <div>
      <div
        classList={{
          'translate-x-0': open(),
          '-translate-x-full': !open()
        }}
        class="fixed left-0 top-0 z-10 h-full w-full transform transition-transform duration-300 ease-in-out"
        onClick={() => setOpen(!open())}
      >
        <div class="flex h-full w-1/2 flex-col justify-start space-y-8 bg-white pl-2 pt-4 md:w-1/4 lg:w-1/6">
          <A href="/training/create">
            <i class="fa-solid fa-person-swimming fa-2xl m-2 text-black"></i>
            <span class="font-bold text-black md:text-xl">Create Training</span>
          </A>
          <A href="/session/create">
            <i class="fa-regular fa-clock fa-2xl m-2 text-black"></i>
            <span class="font-bold text-black md:text-xl">Create Session</span>
          </A>
        </div>
      </div>
      <div
        classList={{ 'bg-black/50': open(), 'pointer-events-none': !open() }}
        class="fixed left-0 top-0 h-full w-screen"
      ></div>
      <div class="mx-auto mt-2 h-full">
        <h1 class="mx-4 text-2xl font-bold">This week's trainings</h1>

        <Show
          when={!details.error}
          fallback={
            <div class="m-4 flex items-center justify-start rounded bg-red-300 p-4 font-bold">
              Couldn't load training details for this week
            </div>
          }
        >
          <Show when={details()?.details?.length === 0}>
            <div class="m-4 flex items-center justify-start rounded bg-blue-200 p-4 font-bold">
              No trainings this week
            </div>
          </Show>

          <For each={details()?.details}>
            {(detail) => (
              <div onClick={() => navigate('/training/' + detail.id)}>
                <DetailCard detail={detail} />
              </div>
            )}
          </For>
          {/* Add space at the bottom, so the buttons dont hide block form */}
          <div class="h-32 w-full"></div>
        </Show>

        <button
          class="fixed bottom-4 right-4 flex h-16 w-16 items-center justify-center rounded-full bg-sky-500 text-white shadow"
          onClick={() => setOpen(!open())}
        >
          <i class="fa-solid fa-pen fa-2xl"></i>
        </button>
      </div>
    </div>
  )
}

export default Home
