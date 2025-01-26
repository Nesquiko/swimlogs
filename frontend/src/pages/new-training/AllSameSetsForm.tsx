import { Component, For, Show } from 'solid-js';
import { SetStoreFunction } from 'solid-js/store';
import {
  EquipmentEnum,
  IntervalTypeEnum,
  LastFinishesTypeEnum,
  NewTrainingSet,
  PauseTypeEnum,
  Start,
} from '~/api/generated';
import { useAppState } from '~/AppContext';
import { EquipmentIcons } from '~/components/Equipment';
import { IconX } from '~/components/icons';
import { Button } from '~/components/ui/button';
import { MobileNumberField } from '~/components/ui/number-field';
import {
  MultiSelect,
  Select,
  SelectContent,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from '~/components/ui/select';
import { SecondsRange } from './common';

type StartType = Start['type'];

interface AllSameSetsFormProps {
  set: NewTrainingSet;
  setSet: SetStoreFunction<NewTrainingSet>;
}

// TODO add the rest of the NewTrainingSet, ended with start
const AllSameSetsForm: Component<AllSameSetsFormProps> = (props) => {
  const { t } = useAppState();

  return (
    <div class="flex flex-col gap-4">
      <MultiSelect<EquipmentEnum>
        values={props.set.equipment ?? []}
        onChange={(eqs) => props.setSet('equipment', eqs)}
        options={Object.values(EquipmentEnum) as EquipmentEnum[]}
        itemComponent={(props) => {
          const EquipmentIcon = EquipmentIcons.get(props.item.rawValue)!;
          return (
            <SelectItem item={props.item}>
              <div class="flex items-center justify-start gap-4">
                <EquipmentIcon
                  size="size-8"
                  class="stroke-1 dark:fill-primary-foreground dark:stroke-primary-foreground"
                />
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
          class="h-fit min-h-14 w-full"
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
                          size="size-8"
                          class="stroke-1 dark:fill-primary-foreground dark:stroke-primary-foreground"
                        />
                        {t(`general.training.equipments.${option}`)}
                        <Button
                          variant="ghost"
                          size="icon"
                          class="size-8"
                          onClick={() => state.remove(option)}
                        >
                          <IconX class="size-5" />
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

      <div class="grid grid-cols-2">
        <Select<StartType>
          value={props.set.start?.type}
          onChange={(typ) => {
            if (!typ) {
              props.setSet('start', undefined);
              return;
            }

            switch (typ) {
              case 'interval':
              case 'pause':
                props.setSet('start', {
                  type: typ,
                  seconds: 0,
                });
                break;
              case 'last-finishes':
                props.setSet('start', {
                  type: LastFinishesTypeEnum.LastFinishes,
                });
            }
          }}
          options={[
            IntervalTypeEnum.Interval,
            PauseTypeEnum.Pause,
            LastFinishesTypeEnum.LastFinishes,
          ]}
          itemComponent={(props) => (
            <SelectItem item={props.item}>
              <div class="flex items-center justify-start gap-2">
                {t(`general.starttype.${props.item.rawValue}`)}
              </div>
            </SelectItem>
          )}
        >
          <SelectLabel class="text-sm">
            {t('general.training.start')}
          </SelectLabel>
          <SelectTrigger aria-label="Fruit" class="w-40">
            <SelectValue<StartType>>
              {(state) => (
                <div class="flex items-center justify-start gap-2">
                  {t(`general.starttype.${state.selectedOption()}`)}
                </div>
              )}
            </SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>

        <Show when={props.set.start} keyed>
          {(start) => {
            if (start.type === LastFinishesTypeEnum.LastFinishes) {
              return;
            }

            return (
              <div>
                <MobileNumberField
                  value={Math.floor(start.seconds / 60)}
                  minValue={SecondsRange.start}
                  maxValue={SecondsRange.end}
                  label={t('general.minutes')}
                  onRawValueChange={(mins) => props.setSet('repeat', mins)}
                  isValid={false}
                />
                <MobileNumberField
                  value={start.seconds}
                  minValue={SecondsRange.start}
                  maxValue={SecondsRange.end}
                  label={t('general.seconds')}
                  onRawValueChange={(secs) => props.setSet('repeat', secs)}
                  isValid={false}
                />
              </div>
            );
          }}
        </Show>
        {/* <IncrementalCounter */}
        {/*   label={t('seconds')} */}
        {/*   value={seconds()} */}
        {/*   onChange={setSeconds} */}
        {/*   min={0} */}
        {/*   max={60} */}
        {/*   step={5} */}
        {/* /> */}
        {/* <IncrementalCounter */}
        {/*   label={t('minutes')} */}
        {/*   value={minutes()} */}
        {/*   onChange={setMinutes} */}
        {/*   min={0} */}
        {/*   max={59} */}
        {/* /> */}
      </div>
    </div>
  );
};

export default AllSameSetsForm;
