import { Trans, useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, For, Show } from 'solid-js'
import IncrementalCounter from '../components/common/IncrementalCounter'
import {
  NumberInput,
  SelectInput,
  TextAreaInput,
} from '../components/common/Input'
import { EquipmentIcons } from '../components/Equipment'
import { EquipmentEnum, NewTrainingSet, StartTypeEnum } from 'swimlogs-api'
import { SmallIntMax } from '../lib/consts'

const DISTANCES = [25, 50, 75, 100, 200, 400]

interface EditSetPageProps {
  set?: NewTrainingSet
  submitLabel: string
  onSubmitSet: (set: NewTrainingSet) => void

  onCancel: () => void
}

const EditSetPage: Component<EditSetPageProps> = (props) => {
  const [t] = useTransContext()
  const [repeat, setRepeat] = createSignal(props.set?.repeat ?? 1)
  const [distance, setDistance] = createSignal(props.set?.distanceMeters ?? 100)
  const distanceErr = () => {
    if (distance() < 25) {
      return 'Too small'
    } else if (distance() > SmallIntMax) {
      return 'Too big'
    }
  }
  const [start, setStart] = createSignal<String>(props.set?.startType ?? 'None')
  const [seconds, setSeconds] = createSignal<number>(
    (props.set?.startSeconds ?? 0) % 60 ?? 0
  )
  const [minutes, setMinutes] = createSignal<number>(
    ((props.set?.startSeconds ?? 0) - seconds()) / 60 ?? 0
  )
  const [equipment, setEquipment] = createSignal<EquipmentEnum[]>(
    props.set?.equipment ?? []
  )
  const [description, setDescription] = createSignal<string | undefined>(
    props.set?.description ?? undefined
  )

  const isValid = () => {
    const isStartValid =
      start() === StartTypeEnum.None ||
      (seconds() === 0 && minutes() > 0) ||
      (seconds() > 0 && minutes() === 0) ||
      (seconds() > 0 && minutes() > 0)

    return distanceErr() === undefined && isStartValid
  }

  const submitSet = () => {
    const startInSeconds = seconds() + minutes() * 60

    const set: NewTrainingSet = {
      setOrder: -1,
      repeat: repeat(),
      distanceMeters: distance(),
      description: description(),
      startType: start() as StartTypeEnum,
      startSeconds: startInSeconds,
      totalDistance: repeat() * distance(),
      equipment: equipment().length > 0 ? equipment() : undefined,
    }

    props.onSubmitSet(set)
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

  const equipmentButton = (eq: EquipmentEnum) => {
    let imgSrc = EquipmentIcons.get(eq)
    return (
      <button
        classList={{
          'bg-sky-300': equipment().includes(eq),
        }}
        class="rounded-lg border border-slate-300 p-2"
        onClick={() => {
          if (equipment().includes(eq)) {
            setEquipment((e) => e.filter((e) => e !== eq))
            return
          }

          setEquipment((e) => [...e, eq])
        }}
      >
        <img width={48} height={48} src={imgSrc} />
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
      <div>
        <SelectInput<String>
          label={t('start')}
          onChange={(opt) => setStart(opt!.value)}
          options={Object.keys(StartTypeEnum).map((st) => {
            return { label: st, value: st }
          })}
        />
        <Show when={start() !== undefined && start() !== StartTypeEnum.None}>
          <div class="space-y-2 pb-2 pl-8">
            <IncrementalCounter
              label={t('seconds')}
              value={seconds()}
              onChange={setSeconds}
              min={0}
              max={60}
              step={5}
            />
            <IncrementalCounter
              label={t('minutes')}
              value={minutes()}
              onChange={setMinutes}
              min={0}
              max={59}
            />
          </div>
        </Show>
        <h1 class="text-xl">{t('equipment')}</h1>
        <div class="flex items-center justify-between">
          {equipmentButton(EquipmentEnum.Fins)}
          {equipmentButton(EquipmentEnum.Monofin)}
          {equipmentButton(EquipmentEnum.Snorkel)}
          {equipmentButton(EquipmentEnum.Paddles)}
          {equipmentButton(EquipmentEnum.Board)}
        </div>
      </div>
      <TextAreaInput
        label={t('description')}
        onInput={(value) => {
          if (value === undefined || value === '') {
            setDescription(undefined)
            return
          }
          setDescription(value)
        }}
        value={description()}
        placeholder={t('set.description.placeholder')}
      />

      <div class="flex items-center justify-between md:justify-around">
        <button
          class="w-24 rounded-lg bg-red-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-red-300"
          onClick={props.onCancel}
        >
          <Trans key="cancel" />
        </button>

        <button
          disabled={!isValid()}
          classList={{
            'bg-green-500/30': !isValid(),
          }}
          class="w-24 rounded-lg bg-green-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-green-300"
          onClick={() => {
            if (!isValid()) {
              return
            }
            submitSet()
          }}
        >
          {props.submitLabel}
        </button>
      </div>
    </div>
  )
}

export default EditSetPage
