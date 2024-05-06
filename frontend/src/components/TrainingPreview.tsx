import { Trans } from '@mbarzda/solid-i18next';
import { type Component, For } from 'solid-js';
import { NewTraining, Training } from 'swimlogs-api';
import SetCard, { Option, SkeletonSetCard } from './SetCard';
import TrainingSummary from './TraningSummary';

interface TrainingPreviewPageProps {
  training: NewTraining | Training;

  setOptions?: Option[];
}

const TrainingPreview: Component<TrainingPreviewPageProps> = (props) => {
  return (
    <div class="space-y-4 px-4">
      <TrainingSummary training={props.training} />

      <h1 class="text-2xl font-bold text-sky-900">
        <Trans key="sets" />
      </h1>
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
