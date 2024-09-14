import { Component } from 'solid-js';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';
import { IconArrowLeft } from '~/components/icons';

// TODO after HomePage do this
const SetFormPage: Component = () => {
  const { t, navigate } = useAppState();
  const [, setTraining] = useNewTraining()!;

  return (
    <div>
      <h1 class="flex items-center gap-2 text-2xl">
        <IconArrowLeft
          class="inline cursor-pointer"
          onClick={() => navigate(-1)}
        />
        <span>New Set</span>
      </h1>
    </div>
  );
};

export default SetFormPage;
