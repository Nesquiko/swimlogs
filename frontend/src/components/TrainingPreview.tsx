import { type Component, For, Show } from 'solid-js';
import { NewTraining, Training } from 'swimlogs-api';
import { locale, minutesToHoursAndMintes } from '../lib/datetime';
import SetCard, { Option, SkeletonSetCard } from './SetCard';
import TrainingSummary from './TraningSummary';

interface TrainingPreviewPageProps {
  training: NewTraining | Training;

  showSession?: boolean;
  setOptions?: Option[];
}

const TrainingPreview: Component<TrainingPreviewPageProps> = (props) => {
  return (
    <div class="space-y-4 px-4">
      <Show when={props.showSession}>
        <div class="grid grid-cols-3">
          <p class="text-xl text-left">
            {props.training.start.toLocaleDateString(locale())}
          </p>
          <p class="text-xl text-center">
            {props.training.start.toLocaleTimeString(locale(), {
              hour: '2-digit',
              minute: '2-digit',
            })}
          </p>
          <p class="text-xl text-right">
            {minutesToHoursAndMintes(props.training.durationMin)}
          </p>
        </div>
      </Show>
      <TrainingSummary training={props.training} />

      <div class="space-y-2">
        <For each={props.training.sets}>
          {(set) => <SetCard set={set} setOptions={props.setOptions} />}
        </For>
      </div>
    </div>
  );
};

const SkeletonTrainingPreview: Component = () => {
  return (
    <div class="space-y-4 px-4 animate-pulse">
      <div class="grid grid-cols-3 items-center">
        <div class="col-start-2 h-8 me-2 inline-block w-full rounded bg-sky-100 px-2.5 py-0.5 text-center text-sky-900"></div>
      </div>
      <div class="space-y-2">
        <For each={Array(8)}>{(_) => <SkeletonSetCard />}</For>
      </div>
    </div>
  );
};

export default TrainingPreview;
export { SkeletonTrainingPreview };
