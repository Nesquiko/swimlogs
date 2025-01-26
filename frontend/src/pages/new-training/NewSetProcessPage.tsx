import { Component, createSignal, Match, Show, Switch } from 'solid-js';
import { createStore, SetStoreFunction } from 'solid-js/store';
import { NewTrainingSet, TypeEnum } from '~/api/generated';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';
import { makePersisted } from '@solid-primitives/storage';
import {
  IconArrowLeft,
  IconArrowRight,
  IconFire,
  IconGroup,
  IconPyramid,
  IconSwimmer,
} from '~/components/icons';
import { DistanceRange, DISTANCES, RepeatRange, REPEATS } from './common';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select';
import OptionsGrid from './new-set-subprocesses/OptionsGrid';
import { CompoundSetForm } from './new-set-subprocesses/compound-set';
import BottomIcons from './new-set-subprocesses/BottomIcons';

const NewSetProcessPage: Component = () => {
  const { t, navigate } = useAppState();
  const [training, setTraining] = useNewTraining()!;
  let [set, setSet] = makePersisted(
    createStore(defaultNewSet(training.sets.length)),
    { name: 'new-set' }
  );

  const [currentStep, setCurrentStep] = createSignal<
    'choosing-type' | TypeEnum
  >('choosing-type');

  const title = () => {
    switch (currentStep()) {
      case 'choosing-type':
      case 'normal':
        return 'newtraining.new.set';
      case 'compound':
        return 'newtraining.new.compound.set';
      case 'pyramid':
        return 'newtraining.new.pyramid.set';
      case 'super':
        return 'newtraining.new.super.set';
    }
  };

  const onSubmitSet = () => {
    setTraining('sets', (sets) => [...sets, set]);
    [set, setSet] = makePersisted(
      createStore(defaultNewSet(training.sets.length)),
      { name: 'new-set' }
    );
  };

  return (
    <div class="h-full">
      <h1 class="flex items-center gap-4 pb-8 text-2xl">
        <IconArrowLeft
          class="inline cursor-pointer"
          onClick={() => navigate(-1)}
        />
        <span>{t(title())}</span>
      </h1>

      <Switch>
        <Match when={currentStep() === 'choosing-type'}>
          <TypeRepeatDistanceForm set={set} setSet={setSet} />
          <BottomIcons
            right={{
              comp: IconArrowRight,
              onClick: () => setCurrentStep(set.type),
            }}
          />
        </Match>

        <Match when={currentStep() === 'compound'}>
          <CompoundSetForm
            set={set}
            setSet={setSet}
            onCancel={() => setCurrentStep('choosing-type')}
            onSubmitSet={onSubmitSet}
          />
        </Match>
      </Switch>
    </div>
  );
};

interface TypeRepeatDistanceFormProps {
  set: NewTrainingSet;
  setSet: SetStoreFunction<NewTrainingSet>;
}

const TypeRepeatDistanceForm: Component<TypeRepeatDistanceFormProps> = (
  props
) => {
  const { t } = useAppState();

  const setTypeIcon = (type: TypeEnum) => {
    switch (type) {
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

  const setTypeDescription = (type: TypeEnum) => {
    let text: string;
    switch (type) {
      case TypeEnum.Normal:
        text =
          'Normal set has a fixed distance. Each repetition can have different intensity, style, start, equipment, ...';
        break;
      case TypeEnum.Compound:
        text =
          'Compound set has a distance which consists of many sub-distances with different intensity, style, start, equipment, ...';
        break;
      case TypeEnum.Pyramid:
        text =
          'Pyramid set involves completing distances that increase up to top distance and then decrease back to the starting one. Each repetition can have different intensity, style, equipment, ...';
        break;
      case TypeEnum.Super:
        text = 'TODO';
        break;
    }
    return <p class="pt-2">{text}</p>;
  };

  const isRepeatValid = () =>
    props.set.repeat >= RepeatRange.start &&
    props.set.repeat <= RepeatRange.end;

  const isDistanceValid = () =>
    props.set.distanceMeters >=
      (DistanceRange.validationStart ?? DistanceRange.start) &&
    props.set.distanceMeters <= DistanceRange.end;

  return (
    <div class="flex flex-col gap-8">
      <div class="h-36">
        <Select<TypeEnum>
          value={props.set.type}
          onChange={(type) => {
            props.setSet('type', type || TypeEnum.Normal);
          }}
          options={Object.values(TypeEnum)}
          itemComponent={(props) => (
            <SelectItem item={props.item}>
              <div class="flex items-center justify-start gap-2">
                {setTypeIcon(props.item.rawValue)}
                {t(`general.training.set.types.${props.item.rawValue}`)}
              </div>
            </SelectItem>
          )}
        >
          <SelectTrigger aria-label="Training type" class="w-full">
            <SelectValue<TypeEnum>>
              {(state) => (
                <div class="flex items-center justify-start gap-2">
                  {setTypeIcon(state.selectedOption())}
                  {t(`general.training.set.types.${state.selectedOption()}`)}
                </div>
              )}
            </SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>
        <p class="text-muted-foreground">
          {setTypeDescription(props.set.type)}
        </p>
      </div>

      <OptionsGrid
        label={t('general.training.repeat')}
        firstRow={[REPEATS[0], REPEATS[1], REPEATS[2]]}
        secondRow={[REPEATS[3], REPEATS[4], REPEATS[5]]}
        value={props.set.repeat}
        onChange={(repeat) => props.setSet('repeat', repeat)}
        orDifferentInput={{
          isValid: isRepeatValid(),
          minValue: RepeatRange.start,
          maxValue: RepeatRange.end,
        }}
      />

      <OptionsGrid
        label={t('newtraining.total.distance')}
        firstRow={[DISTANCES[0], DISTANCES[1], DISTANCES[2]]}
        secondRow={[DISTANCES[3], DISTANCES[4], DISTANCES[5]]}
        value={props.set.distanceMeters}
        onChange={(dist) => props.setSet('distanceMeters', dist)}
        orDifferentInput={{
          isValid: isDistanceValid(),
          minValue: DistanceRange.start,
          maxValue: DistanceRange.end,
          step: 25,
        }}
      />
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

export default NewSetProcessPage;
