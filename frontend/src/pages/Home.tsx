import { Component } from 'solid-js'

const Home: Component = () => {
  return (
    <div>
      <h1>Home</h1>
      <div class="absolute bottom-4 right-4 flex h-16 w-16 items-center justify-center rounded-lg bg-sky-500 shadow">
        <div class="text-3xl text-white">
          <i class="fa-solid fa-pen-to-square"></i>
        </div>
      </div>
    </div>
  )
}

export default Home
