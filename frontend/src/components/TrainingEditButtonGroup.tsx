import { Component, Show } from 'solid-js';
import { NewTraining, Training } from 'swimlogs-api';

interface TrainingEditButtonGroupProps {
  training: NewTraining | Training;

  backLabel: string;
  onBack: () => void;

  onAddSet: () => void;

  confirmLabel: string;
  onConfirm: () => void;
}

const TrainingEditButtonGroup: Component<TrainingEditButtonGroupProps> = (
  props
) => {
  return (
    <div
      class="flex items-center p-4"
      classList={{
        'justify-between md:justify-around': props.training.sets.length > 0,
        'justify-center': props.training.sets.length === 0,
      }}
    >
      <Show when={props.training.sets.length > 0}>
        <button
          class="w-24 rounded-lg bg-red-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-red-300"
          onClick={props.onBack}
        >
          {props.backLabel}
        </button>
      </Show>
      <div class="py-2 text-center">
        <button
          class="h-12 w-12 rounded-full bg-sky-500"
          onClick={props.onAddSet}
        >
          <i class="fa-solid fa-plus fa-2xl text-white"></i>
        </button>
      </div>

      <Show when={props.training.sets.length > 0}>
        <button
          class="w-24 rounded-lg bg-green-500 py-2 text-xl font-bold text-white focus:outline-none focus:ring-2 focus:ring-green-300"
          onClick={props.onConfirm}
        >
          {props.confirmLabel}
        </button>
      </Show>
    </div>
  );
};
export default TrainingEditButtonGroup;
