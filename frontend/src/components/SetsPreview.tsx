import { type Component, For } from 'solid-js';
import { NewTrainingSet, TrainingSet } from 'swimlogs-api';
import SetCard, { Option, SkeletonSetCard } from './SetCard';

interface SetsPreviewProps<T extends NewTrainingSet | TrainingSet> {
  sets: T[];
  setOptions?: Option[];
}

const SetsPreview = <T extends NewTrainingSet | TrainingSet>(
  props: SetsPreviewProps<T>
) => {
  return (
    <div class="space-y-2">
      <For each={props.sets}>
        {(set) => <SetCard set={set} setOptions={props.setOptions} />}
      </For>
    </div>
  );
};

const SkeletonSetsPreview: Component = () => {
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

export default SetsPreview;
export { SkeletonSetsPreview };
