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
    <div class="space-y-2 animate-pulse">
      <For each={Array(8)}>{(_) => <SkeletonSetCard />}</For>
    </div>
  );
};

export default SetsPreview;
export { SkeletonSetsPreview };
