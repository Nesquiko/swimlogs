import { Trans, useTransContext } from '@mbarzda/solid-i18next';
import { Component, For } from 'solid-js';
import InlineDatepicker from '../components/InlineDatepicker';
import { NewTraining } from 'swimlogs-api';
import { SetStoreFunction } from 'solid-js/store';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from './ui/select';

export const StartTimeHours = [
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
];

export const StartTimeMinutes = ['00', '15', '30', '45'];

type Time = { label: string; value: number };

const TIMES = [
  { label: '30 min', value: 30 },
  { label: '1 h', value: 60 },
  { label: '1 h 30 min', value: 90 },
  { label: '2 h', value: 120 },
] as Time[];

interface SessionEditFormProps {
  training: NewTraining;
  updateTraining: SetStoreFunction<NewTraining>;
}

const SessionEditForm: Component<SessionEditFormProps> = ({
  training,
  updateTraining,
}) => {
  const [t] = useTransContext();

  const timeItem = (time: Time) => {
    return (
      <button
        classList={{ 'bg-sky-300': training.durationMin === time.value }}
        class="w-32 rounded-lg border border-slate-300 p-2 text-center text-lg shadow"
        onClick={() => updateTraining('durationMin', time.value)}
      >
        {time.label}
      </button>
    );
  };

  return (
    <div class="space-y-4">
      <h1 class="py-2 text-2xl font-bold text-sky-900">
        <Trans key="training.details" />
      </h1>
      <div>
        <h1 class="text-xl">{t('starttime')}</h1>
        <div class="flex justify-around">
          <Select
            value={training.start.getHours().toString().padStart(2, '0')}
            onChange={(h) => {
              const d = new Date(training.start);
              d.setHours(parseInt(h));
              updateTraining('start', d);
            }}
            options={StartTimeHours}
            itemComponent={(props) => (
              <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
            )}
          >
            <SelectTrigger aria-label="hour" class="w-32 text-xl">
              <SelectValue<string>>
                {(state) => state.selectedOption()}
              </SelectValue>
            </SelectTrigger>
            <SelectContent />
          </Select>

          <Select
            value={training.start.getMinutes().toString().padStart(2, '0')}
            onChange={(m) => {
              const d = new Date(training.start);
              d.setMinutes(parseInt(m));
              updateTraining('start', d);
            }}
            options={StartTimeMinutes}
            itemComponent={(props) => (
              <SelectItem item={props.item}>{props.item.rawValue}</SelectItem>
            )}
          >
            <SelectTrigger aria-label="minute" class="w-32 text-xl">
              <SelectValue<string>>
                {(state) => state.selectedOption()}
              </SelectValue>
            </SelectTrigger>
            <SelectContent />
          </Select>
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
        <InlineDatepicker
          initialDate={training.start}
          onChange={(d) => {
            const newStart = new Date(training.start);
            newStart.setMonth(d.getMonth());
            newStart.setDate(d.getDate());
            updateTraining('start', newStart);
          }}
        />
      </div>
    </div>
  );
};

export default SessionEditForm;
