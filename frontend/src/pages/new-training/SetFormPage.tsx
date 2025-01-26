import { Component, Match, Switch } from 'solid-js';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';
import {
  IconArrowLeft,
  IconFire,
  IconGroup,
  IconPyramid,
  IconSwimmer,
} from '~/components/icons';
import { NewTrainingSet, TypeEnum } from '~/api/generated';
import { createStore } from 'solid-js/store';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select';
import NormalSetForm from './NormalSetForm';

const SetFormPage: Component = () => {
  const { t, navigate } = useAppState();
  const [training, setTraining] = useNewTraining()!;
  const [set, setSet] = createStore<NewTrainingSet>(
    defaultNewSet(training.sets.length)
  );

  return (
    <div class="flex flex-col gap-4">
      <div class="flex justify-between">
        <h1 class="flex items-center gap-2 pb-8 text-2xl">
          <IconArrowLeft
            class="inline cursor-pointer"
            onClick={() => navigate(-1)}
          />
          <span>{t('newtraining.new.set')}</span>
        </h1>
        <Select<TypeEnum>
          value={set.type}
          onChange={(typ) => setSet('type', typ || TypeEnum.Normal)}
          options={Object.values(TypeEnum)}
          itemComponent={(props) => (
            <SelectItem item={props.item}>
              <div class="flex items-center justify-start gap-2">
                <SetTypeIcon typ={props.item.rawValue} />
                {t(`general.training.set.types.${props.item.rawValue}`)}
              </div>
            </SelectItem>
          )}
        >
          <SelectTrigger aria-label="Fruit" class="w-[180px]">
            <SelectValue<TypeEnum>>
              {(state) => (
                <div class="flex items-center justify-start gap-2">
                  <SetTypeIcon typ={state.selectedOption()} />
                  {t(`general.training.set.types.${state.selectedOption()}`)}
                </div>
              )}
            </SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>
      </div>

      <Switch>
        <Match when={set.type === TypeEnum.Normal}>
          <NormalSetForm set={set} setSet={setSet} />
        </Match>
      </Switch>

      <pre>{JSON.stringify(set, null, 2)}</pre>
    </div>
  );
};

const defaultNewSet = (setOrder: number): NewTrainingSet => {
  return {
    setOrder,
    type: TypeEnum.Normal,
    repeat: 1,
    distanceMeters: 100,
    equipment: [],
  };
};

const SetTypeIcon: Component<{ typ: TypeEnum }> = (props) => {
  switch (props.typ) {
    case TypeEnum.Normal:
      return <IconSwimmer />;
    case TypeEnum.Compound:
      return <IconGroup />;
    case TypeEnum.Pyramid:
      return <IconPyramid />;
    case TypeEnum.Super:
      return <IconFire />;
  }
};

export default SetFormPage;
