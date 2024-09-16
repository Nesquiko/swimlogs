import { Component, For } from 'solid-js';
import { useAppState } from '~/AppContext';
import { useNewTraining } from './NewTrainingContext';
import { IconArrowLeft, IconX } from '~/components/icons';
import { EquipmentEnum, NewTrainingSet, TypeEnum } from '~/api/generated';
import { createStore } from 'solid-js/store';
import { MobileNumberField } from '~/components/ui/number-field';
import { Button } from '~/components/ui/button';
import { Label } from '~/components/ui/label';
import { randomId } from '~/lib/str';
import { EquipmentIcons } from '~/components/Equipment';
import {
  MultiSelect,
  SelectContent,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select';

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

      <MultiSelect<EquipmentEnum>
        values={set.equipment ?? []}
        onChange={(eqs) => {
          setSet('equipment', eqs);
        }}
        options={Object.values(EquipmentEnum) as EquipmentEnum[]}
        itemComponent={(props) => {
          const EquipmentIcon = EquipmentIcons.get(props.item.rawValue)!;
          return (
            <SelectItem item={props.item}>
              <div class="flex items-center justify-start gap-4">
                <EquipmentIcon class="stroke-1 dark:fill-primary-foreground dark:stroke-primary-foreground" />
                <span>
                  {/* @ts-ignore option values are taken from EquipmentEnum */}
                  {t(`general.training.equipments.${props.item.textValue}`)}
                </span>
              </div>
            </SelectItem>
          );
        }}
      >
        <SelectLabel class="text-sm">
          {t('general.training.equipment')}
        </SelectLabel>
        <SelectTrigger
          aria-label="equipment"
          as="div"
          class="h-fit min-h-16 w-full"
        >
          <SelectValue<EquipmentEnum>>
            {(state) => (
              <div>
                <For each={state.selectedOptions()}>
                  {(option) => {
                    const EquipmentIcon = EquipmentIcons.get(option)!;
                    return (
                      <span
                        class="inline-flex items-center justify-start"
                        onPointerDown={(e) => e.stopPropagation()}
                      >
                        <EquipmentIcon
                          size="size-10"
                          class="stroke-1 dark:fill-primary-foreground dark:stroke-primary-foreground"
                        />
                        {t(`general.training.equipments.${option}`)}
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => state.remove(option)}
                        >
                          <IconX size={20} />
                        </Button>
                      </span>
                    );
                  }}
                </For>
              </div>
            )}
          </SelectValue>
        </SelectTrigger>
        <SelectContent />
      </MultiSelect>

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
    equipment: [
      EquipmentEnum.Paddles,
      EquipmentEnum.Fins,
      EquipmentEnum.Board,
      EquipmentEnum.Monofin,
      EquipmentEnum.Parachute,
      EquipmentEnum.PullBuoy,
      EquipmentEnum.Snorkel,
    ],
  };
};

export default SetFormPage;
