import { useTransContext } from '@mbarzda/solid-i18next';
import { Component, For, Show } from 'solid-js';
import { NewTrainingSet, StartTypeEnum, TrainingSet } from 'swimlogs-api';
import DropdownMenu from './DropdownMenu';
import { EquipmentIcons } from '../lib/equipment-svgs';
import { GroupColors } from '../lib/set-groups';

interface Option {
  text: string;
  icon: string;
  onClick: (setIdx: number) => void;
  disabled?: boolean;
  disabledFunc?: (setIdx: number) => boolean;
}

interface SetCardProps {
  set: NewTrainingSet | TrainingSet;
  setOptions?: Option[];
}
const SetCard: Component<SetCardProps> = (props) => {
  const [t] = useTransContext();

  const setContent =
    props.set.repeat > 1
      ? `${props.set.repeat}x${props.set.distanceMeters}m`
      : `${props.set.distanceMeters}m`;

  const start = (): string => {
    if (props.set.startSeconds === undefined) {
      return '';
    }

    const startSeconds =
      props.set.startType !== StartTypeEnum.None
        ? props.set.startSeconds % 60
        : 0;
    const startMinutes =
      props.set.startType !== StartTypeEnum.None
        ? (props.set.startSeconds - startSeconds) / 60
        : 0;
    return `${t(props.set.startType.toLowerCase())}: ${
      startMinutes !== 0 ? startMinutes + "'" : ''
    }${startSeconds !== 0 ? startSeconds + '"' : ''}`;
  };

  return (
    <div class="mx-auto block max-w-xl rounded-lg border border-gray-200 bg-white shadow">
      <div
        classList={{
          'rounded-lg':
            props.set.description === undefined &&
            (props.set.equipment === undefined ||
              props.set.equipment.length === 0),
        }}
        class="w-full rounded-t-lg bg-sky-200 p-2 text-sky-900"
      >
        <div class="grid grid-cols-4 gap-x-2 items-center">
          <span class="text-xl font-bold col-start-1">{setContent}</span>
          <Show when={props.set.startType !== StartTypeEnum.None}>
            <span class="text-center col-start-2 col-span-2 text-xl font-bold">
              {start()}
            </span>
          </Show>
          <Show when={props.set.group} keyed>
            {(group) => (
              <span
                classList={{ [GroupColors.get(group)!]: true }}
                class="col-start-4 rounded-lg text-white bg-sky-500 text-center w-full h-6"
              >
                {t(group.toLowerCase())}
              </span>
            )}
          </Show>

          <Show
            when={props.setOptions !== undefined && props.setOptions.length > 0}
          >
            <div class="col-start-5">
              <DropdownMenu
                icon="fa-ellipsis"
                items={(props.setOptions ?? []).slice(0, -1).map((option) => ({
                  icon: option.icon,
                  text: option.text,
                  onClick: () => option.onClick(props.set.setOrder!),
                  disabled:
                    option.disabled ||
                    option.disabledFunc?.(props.set.setOrder!) ||
                    false,
                }))}
                finalItem={{
                  icon: props.setOptions![props.setOptions!.length - 1].icon,
                  text: props.setOptions![props.setOptions!.length - 1].text,
                  onClick: () =>
                    props.setOptions![props.setOptions!.length - 1].onClick(
                      props.set.setOrder!
                    ),
                  disabled:
                    props.setOptions![props.setOptions!.length - 1].disabled ||
                    props.setOptions![
                      props.setOptions!.length - 1
                    ].disabledFunc?.(props.set.setOrder!) ||
                    false,
                }}
              />
            </div>
          </Show>
        </div>
      </div>
      <Show when={props.set.description}>
        <p class="whitespace-pre-wrap p-2 text-gray-500">
          {props.set.description}
        </p>
      </Show>
      <Show when={props.set.equipment && props.set.equipment.length > 0}>
        <div class="text-center">
          <For each={props.set.equipment}>
            {(equipment) => (
              <img
                class="inline-block"
                width={64}
                height={64}
                src={EquipmentIcons.get(equipment)}
              />
            )}
          </For>
        </div>
      </Show>
    </div>
  );
};

const SkeletonSetCard: Component = () => {
  const equipmentSkeleton = () => {
    return <i class="fa-solid fa-image fa-2xl inline-block text-gray-200" />;
  };

  return (
    <div class="mx-auto block max-w-xl rounded-lg border border-gray-200 bg-white shadow animate-pulse">
      <div class="w-full rounded-t-lg bg-sky-100 p-2 text-sky-900">
        <div class="h-6 space-x-6">
          <span class="inline-block h-2 bg-sky-200 rounded w-20"></span>
          <span class="inline-block h-2 bg-sky-200 rounded w-28"></span>
        </div>
      </div>
      <div class="space-y-3 p-2">
        <div class="h-2 bg-gray-200 rounded"></div>
        <div class="grid grid-cols-3 gap-4">
          <div class="h-2 bg-gray-200 rounded col-span-2"></div>
          <div class="h-2 bg-gray-200 rounded col-span-1"></div>
        </div>
      </div>

      <div class="text-center space-x-4 p-2">
        <For each={Array(3)}>{equipmentSkeleton}</For>
      </div>
    </div>
  );
};

export default SetCard;
export type { Option };
export { SkeletonSetCard };
