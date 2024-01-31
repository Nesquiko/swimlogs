import { useTransContext } from '@mbarzda/solid-i18next'
import { Component, createSignal, For } from 'solid-js'
import { SelectInput } from '../components/common/Input'
import InlineDatepicker from '../components/InlineDatepicker'

const StartTimeHours = [
  '06',
  '07',
  '08',
  '09',
  '10',
  '11',
  '12',
  '13',
  '14',
  '15',
  '16',
  '17',
  '18',
  '19',
  '20',
]

const StartTimeMinutes = ['00', '15', '30', '45']

type Time = { label: string; value: number }

const TIMES = [
  { label: '30 min', value: 30 },
  { label: '1 h', value: 60 },
  { label: '1 h 30 min', value: 90 },
  { label: '2 h', value: 120 },
] as Time[]

interface EditTrainingSessionPageProps {
  initial?: { start: Date; durationMin: number }
  onSubmit: (session: { start: Date; durationMin: number }) => void
  onBack: () => void
}

const EditTrainingSessionPage: Component<EditTrainingSessionPageProps> = (
  props
) => {
  const [t] = useTransContext()
  const [durationMin, setDurationMin] = createSignal(
    props.initial?.durationMin || 60
  )
  const [date, setDate] = createSignal(props.initial?.start || new Date())
  const [hours, setHours] = createSignal<String>(
    props.initial?.start.getHours().toString().padStart(2, '0') || '18'
  )
  const [minutes, setMinutes] = createSignal<String>(
    props.initial?.start.getMinutes().toString().padStart(2, '0') || '00'
  )

  const timeItem = (time: Time) => {
    return (
      <button
        classList={{ 'bg-sky-300': durationMin() === time.value }}
        class="w-32 rounded-lg border border-slate-300 p-2 text-center text-lg shadow"
        onClick={() => setDurationMin(time.value)}
      >
        {time.label}
      </button>
    )
  }

  const onSubmit = () => {
    const start = new Date(date())
    start.setHours(Number(hours()))
    start.setMinutes(Number(minutes()))
    start.setSeconds(0)

    props.onSubmit({
      start,
      durationMin: durationMin(),
    })
  }

  return (
    <div class="space-y-4 px-4">
      <div>
        <h1 class="text-xl">{t('starttime')}</h1>
        <div class="flex justify-around">
          <SelectInput<String>
            onChange={(h) => setHours(h!.value)}
            initialValueIndex={StartTimeHours.findIndex((h) => h === hours())}
            options={Array.from(StartTimeHours).map((h) => ({
              label: h,
              value: h,
            }))}
          />

          <SelectInput<String>
            onChange={(m) => setMinutes(m!.value)}
            initialValueIndex={StartTimeMinutes.findIndex(
              (m) => m === minutes()
            )}
            options={Array.from(StartTimeMinutes).map((m) => ({
              label: m,
              value: m,
            }))}
          />
        </div>
      </div>

      <div>
        <h1 class="text-xl">{t('duration')}</h1>
        <div class="grid grid-cols-2 justify-items-center gap-4 md:grid-cols-4">
          <For each={TIMES}>{timeItem}</For>
        </div>
      </div>

      <div>
        <h1 class="text-xl">{t('date')}</h1>
        <InlineDatepicker initialDate={date()} onChange={setDate} />
      </div>

      <div class="flex items-center justify-between md:justify-around">
        <button
          class="w-24 rounded-lg bg-red-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-red-300"
          onClick={props.onBack}
        >
          {t('back')}
        </button>

        <button
          class="w-24 rounded-lg bg-green-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-green-300"
          onClick={onSubmit}
        >
          {t('finish')}
        </button>
      </div>
    </div>
  )
}

export default EditTrainingSessionPage
