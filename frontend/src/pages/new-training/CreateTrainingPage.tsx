import { type Component } from 'solid-js';
import { Button } from '~/components/ui/button';
import { IconArrowLeft, IconPlus } from '~/components/icons';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';

const CreateTrainingPage: Component = () => {
  const { t, navigate, setGoBackTo } = useAppState();
  const [training, setTraining] = useNewTraining()!;

  return (
    <div>
      <h1 class="flex items-center gap-2 text-2xl">
        <IconArrowLeft
          class="inline cursor-pointer"
          onClick={() => navigate(-1)}
        />
        <span>{t('newtraining.new.training')}</span>
      </h1>
      <Button
        size="icon"
        class="fixed bottom-4 right-4"
        onClick={() => navigate('/training/create/set/new')}
      >
        <IconPlus size={30} />
      </Button>
    </div>
  );
};

export default CreateTrainingPage;
