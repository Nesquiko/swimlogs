import { useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, For } from 'solid-js'
import IncrementalCounter from '../components/common/IncrementalCounter'
import { NumberInput } from '../components/common/Input'
import { SmallIntMax } from '../lib/consts'

const DISTANCES = [25, 50, 75, 100, 200, 400]

const CreateSetPage: Component = () => {
  const [t] = useTransContext()
  const [repeat, setRepeat] = createSignal(1)
  const [distance, setDistance] = createSignal(100)

  const distanceErr = () => {
    if (distance() < 25) {
      return 'Too small'
    } else if (distance() > SmallIntMax) {
      return 'Too big'
    }
  }

  const distanceItem = (dist: number) => {
    return (
      <button
        classList={{ 'bg-sky-300': distance() === dist }}
        class="w-24 rounded-lg border border-slate-300 p-2 text-center text-lg shadow"
        onClick={() => setDistance(dist)}
      >
        {dist}
      </button>
    )
  }

  return (
    <div class="space-y-4 px-4">
      <IncrementalCounter
        label={t('repeat')}
        value={repeat()}
        onChange={setRepeat}
        min={1}
        max={SmallIntMax}
      />
      <h1 class="text-xl">{t('distance')}</h1>
      <div class="space-y-4">
        <div class="grid grid-cols-[100px,100px,100px] justify-between justify-items-center gap-4 text-center">
          <For each={DISTANCES.slice(0, 3)}>{distanceItem}</For>
          <For each={DISTANCES.slice(3)}>{distanceItem}</For>
        </div>
        <NumberInput
          placeholder={t('different.distance')}
          onChange={(n) => {
            if (n === undefined) {
              n = 25
            }
            setDistance(n)
          }}
          value={DISTANCES.includes(distance()) ? undefined : distance()}
          error={distanceErr()}
        />
      </div>
      <h1 class="text-xl">{t('start')}</h1>

      <pre>
        {JSON.stringify({
          repeat: repeat(),
          distance: distance(),
          distanceErr,
        })}
      </pre>
    </div>
  )
}

export default CreateSetPage
