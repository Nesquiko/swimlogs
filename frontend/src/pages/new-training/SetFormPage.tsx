import { Accessor, Component, For } from 'solid-js';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';
import { IconArrowLeft } from '~/components/icons';
import { EquipmentEnum, NewTrainingSet, TypeEnum } from '~/api/generated';
import { createStore } from 'solid-js/store';
import { MobileNumberField } from '~/components/ui/number-field';
import { Button } from '~/components/ui/button';
import { Label } from '~/components/ui/label';
import { randomId } from '~/lib/str';
import { EquipmentButton } from '~/components/Equipment';

const DISTANCES = [25, 50, 75, 100, 200, 400];

interface ValidationRange {
  start: number;
  validationStart?: number;
  end: number;
}

const RepeatRange: ValidationRange = { start: 1, end: 500 };
const DistanceRange: ValidationRange = {
  start: 0,
  validationStart: 1,
  end: 30000,
};

const SetFormPage: Component = () => {
  const { t, navigate } = useAppState();
  const [training, setTraining] = useNewTraining()!;
  const [set, setSet] = createStore<NewTrainingSet>(
    defaultNewSet(training.sets.length)
  );

  const isRepeatValid = () =>
    set.repeat >= RepeatRange.start && set.repeat <= RepeatRange.end;
  const isDistanceValid = () =>
    set.distanceMeters >=
      (DistanceRange.validationStart ?? DistanceRange.start) &&
    set.distanceMeters <= DistanceRange.end;

  const distanceItem = (dist: number) => {
    return (
      <Button
        variant={set.distanceMeters === dist ? 'outline' : 'ghost'}
        onClick={() => setSet('distanceMeters', dist)}
      >
        {dist}
      </Button>
    );
  };
  const distanceId = randomId();

  const equipmentButton = (eq: EquipmentEnum) => {
    return (
      <EquipmentButton
        eq={eq}
        isActive={set.equipment?.includes(eq) ?? false}
        onClick={() => {
          if (!set.equipment) {
            setSet('equipment', [eq]);
            return;
          }

          if (set.equipment.includes(eq)) {
            setSet('equipment', (eqs) => eqs?.filter((e) => e !== eq));
            return;
          }

          setSet('equipment', [...set.equipment, eq]);
        }}
      />
    );
  };

  return (
    <div class="flex flex-col gap-4">
      <h1 class="flex items-center gap-2 text-2xl">
        <IconArrowLeft
          class="inline cursor-pointer"
          onClick={() => navigate(-1)}
        />
        <span>{t('newtraining.new.set')}</span>
      </h1>

      <div class="grid grid-cols-2">
        <MobileNumberField
          value={set.repeat}
          minValue={RepeatRange.start}
          maxValue={RepeatRange.end}
          label={t('general.training.repeat')}
          onRawValueChange={(repeat) => setSet('repeat', repeat)}
          isValid={isRepeatValid()}
        />
        <div class="grid grid-cols-3">
          <Label class="col-span-3 pb-4" for={`${distanceId}-input`}>
            {t('general.training.distance')}
          </Label>
          <For each={DISTANCES.slice(0, 3)}>{distanceItem}</For>
          <For each={DISTANCES.slice(3)}>{distanceItem}</For>
          <div class="relative col-span-3 flex items-center py-1 text-sm">
            <div class="flex-grow border-t"></div>
            <span class="mx-2 text-muted-foreground">
              {t('newtraining.or.different')}
            </span>
            <div class="flex-grow border-t"></div>
          </div>
          <MobileNumberField
            id={distanceId}
            class="col-span-3 w-full"
            value={set.distanceMeters}
            minValue={DistanceRange.start}
            maxValue={DistanceRange.end}
            step={25}
            onRawValueChange={(dist) => setSet('distanceMeters', dist)}
            isValid={isDistanceValid()}
          />
        </div>
      </div>

      <Label>{t('general.training.equipment')}</Label>
      <div class="grid grid-cols-3 gap-4">
        {equipmentButton(EquipmentEnum.Snorkel)}
        {equipmentButton(EquipmentEnum.Board)}
        {equipmentButton(EquipmentEnum.Paddles)}
      </div>
      <div class="grid grid-cols-4 gap-4">
        {equipmentButton(EquipmentEnum.Fins)}
        {equipmentButton(EquipmentEnum.Monofin)}
        {equipmentButton(EquipmentEnum.PullBuoy)}
        {equipmentButton(EquipmentEnum.Parachute)}
      </div>

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

export default SetFormPage;
