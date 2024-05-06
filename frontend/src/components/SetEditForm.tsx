import { useTransContext } from '@mbarzda/solid-i18next';
import { Component, For, Show } from 'solid-js';
import IncrementalCounter from '../components/IncrementalCounter';
import { NumberInput, SelectInput, TextAreaInput } from '../components/Input';
import { EquipmentIcons } from '../lib/equipment-svgs';
import {
  EquipmentEnum,
  NewTrainingSet,
  StartTypeEnum,
  GroupEnum,
} from 'swimlogs-api';
import { SmallIntMax } from '../lib/consts';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '../components/ui/select';
import { Label } from '../components/ui/label';
import { SetStoreFunction } from 'solid-js/store';

const DISTANCES = [25, 50, 75, 100, 200, 400];
const NO_GROUP_KEY = 'no.group';

interface SetEditFormProps {
  set: NewTrainingSet;
  updateSet: SetStoreFunction<NewTrainingSet>;
}

const SetEditForm: Component<SetEditFormProps> = ({ set, updateSet }) => {
  const [t] = useTransContext();
  const distanceErr = () => {
    if (set.distanceMeters < 25) {
      return t('distance.error.too.small');
    } else if (set.distanceMeters > SmallIntMax) {
      return t('distance.error.too.big');
    }
  };

  const distanceItem = (dist: number) => {
    return (
      <button
        classList={{ 'bg-sky-300': set.distanceMeters === dist }}
        class="w-24 rounded-lg border border-slate-300 p-2 text-center text-lg shadow"
        onClick={() => updateSet('distanceMeters', dist)}
      >
        {dist}
      </button>
    );
  };

  const equipmentButton = (eq: EquipmentEnum) => {
    let imgSrc = EquipmentIcons.get(eq);
    return (
      <button
        classList={{
          'bg-sky-300': set.equipment?.includes(eq),
        }}
        class="rounded-lg border border-slate-300 p-2"
        onClick={() => {
          if (set.equipment?.includes(eq)) {
            updateSet('equipment', (equipment) =>
              equipment?.filter((e) => e !== eq)
            );
            return;
          }

          updateSet('equipment', (equipment) => [...(equipment ?? []), eq]);
        }}
      >
        <img width={48} height={48} src={imgSrc} />
      </button>
    );
  };

  return (
    <div class="space-y-4">
      <IncrementalCounter
        label={t('repeat')}
        value={set.repeat}
        onChange={(repeat) => updateSet('repeat', repeat)}
        min={1}
        max={SmallIntMax}
      />
      <h1 class="text-xl">{t('distance')}</h1>
      <div class="space-y-4">
        <div class="grid grid-cols-[100px,100px,100px] justify-between justify-items-center gap-4 text-center">
          <For each={DISTANCES.slice(0, 3)}>{distanceItem}</For>
          <For each={DISTANCES.slice(3)}>{distanceItem}</For>
        </div>
        <NumberInput
          placeholder={t('different.distance')}
          onChange={(n) => {
            if (n === undefined) {
              n = 25;
            }
            updateSet('distanceMeters', n);
          }}
          value={
            DISTANCES.includes(set.distanceMeters)
              ? undefined
              : set.distanceMeters
          }
          error={distanceErr()}
        />
      </div>
      <SelectInput<String>
        label={t('start')}
        onChange={(opt) => updateSet('startType', opt!.value as StartTypeEnum)}
        initialValueIndex={Object.keys(StartTypeEnum).findIndex(
          (s) => s === set.startType
        )}
        options={Object.keys(StartTypeEnum).map((st) => {
          return { label: t(st.toLowerCase()), value: st };
        })}
      />
      <Show when={set.startType !== StartTypeEnum.None}>
        <div class="space-y-2 pb-2 pl-8">
          <IncrementalCounter
            label={t('seconds')}
            value={(set.startSeconds ?? 0) % 60}
            onChange={(s) =>
              updateSet(
                'startSeconds',
                (secs) => Math.floor((secs ?? 0) / 60) * 60 + s
              )
            }
            min={0}
            max={55}
            step={5}
          />
          <IncrementalCounter
            label={t('minutes')}
            value={Math.floor((set.startSeconds ?? 0) / 60)}
            onChange={(m) =>
              updateSet('startSeconds', (secs) => ((secs ?? 0) % 60) + m * 60)
            }
            min={0}
            max={59}
          />
        </div>
      </Show>
      <h1 class="text-xl">{t('equipment')}</h1>
      <div class="flex items-center justify-between">
        {equipmentButton(EquipmentEnum.Fins)}
        {equipmentButton(EquipmentEnum.Monofin)}
        {equipmentButton(EquipmentEnum.Snorkel)}
        {equipmentButton(EquipmentEnum.Paddles)}
        {equipmentButton(EquipmentEnum.Board)}
      </div>
      <TextAreaInput
        label={t('description')}
        onInput={(value) => {
          if (value === undefined || value === '') {
            updateSet('description', undefined);
            return;
          }
          updateSet('description', value);
        }}
        value={set.description}
        placeholder={t('set.description.placeholder')}
      />

      <div class="w-full flex justify-between items-center">
        <Label for="group" class="text-xl">
          {t('group')}
        </Label>
        <Select
          id="group"
          value={set.group ?? (NO_GROUP_KEY as GroupEnum)}
          onChange={(g) => {
            if (g === (NO_GROUP_KEY as GroupEnum) || g === null) {
              updateSet('group', undefined);
              return;
            }
            updateSet('group', g);
          }}
          options={[NO_GROUP_KEY].concat(...Object.values(GroupEnum))}
          placeholder={
            <span class="text-black/50">{t('group.placeholder')}</span>
          }
          itemComponent={(props) => (
            <SelectItem class="text-lg" item={props.item}>
              {t(props.item.rawValue)}
            </SelectItem>
          )}
        >
          <SelectTrigger aria-label="group" class="w-44 text-lg">
            <SelectValue<GroupEnum>>{(g) => t(g.selectedOption())}</SelectValue>
          </SelectTrigger>
          <SelectContent />
        </Select>
      </div>
    </div>
  );
};

export default SetEditForm;
