import { Trans, useTransContext } from '@mbarzda/solid-i18next';
import { Component, For, Show } from 'solid-js';
import {
  GroupEnum,
  NewTraining,
  NewTrainingSet,
  Training,
  TrainingSet,
} from 'swimlogs-api';
import { locale, minutesToHoursAndMintes } from '../lib/datetime';
import { GroupColors } from '../lib/set-groups';
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from './ui/card';

interface TrainingSummaryProps {
  training: NewTraining | Training;
  onClick?: () => void;
}

type DistancePerGroup = { group: GroupEnum; distance: number };

const TrainingSummary: Component<TrainingSummaryProps> = (props) => {
  const [t] = useTransContext();
  const distances = distancePerGroup(props.training.sets);

  const uniqueGroups = props.training.sets
    .map((s) => s.group)
    .filter((g) => g).length;

  const labelKey =
    uniqueGroups === 0 ? 'distance.in.training' : 'distance.per.group';

  const distanceItem = (dist: DistancePerGroup) => {
    const color = GroupColors.get(dist.group as GroupEnum)!;

    return (
      <div>
        <span
          classList={{ [color]: true }}
          class="inline-block rounded-lg text-center w-1/2 h-6 text-white"
        >
          {t(dist.group.toLowerCase())}
        </span>

        <span class="inline-block text-center w-1/2 h-6 text-black">
          {dist.distance}m
        </span>
      </div>
    );
  };

  return (
    <Card
      classList={{ 'cursor-pointer': props.onClick !== undefined }}
      class="max-w-lg"
      onClick={props.onClick}
    >
      <CardHeader>
        <CardTitle class="text-xl text-sky-900 flex justify-between">
          <span>
            {props.training.start
              .toLocaleDateString(locale())
              .replaceAll(' ', '')}{' '}
            <Trans key="at" />{' '}
            {props.training.start.toLocaleTimeString(locale(), {
              hour: '2-digit',
              minute: '2-digit',
            })}
          </span>
          <span class="font-normal space-x-1">
            <i class="fa-solid fa-clock" />
            <b>{minutesToHoursAndMintes(props.training.durationMin)}</b>
          </span>
        </CardTitle>
      </CardHeader>

      <Show
        when={uniqueGroups !== 0}
        fallback={
          <CardFooter class="text-lg flex justify-between">
            <Trans key={labelKey} />
            <b class="inline-block">{props.training.totalDistance}m</b>
          </CardFooter>
        }
      >
        <CardContent class="space-y-2">
          <span class="text-lg">
            <Trans key={labelKey} />
          </span>
          <For each={distances}>{distanceItem}</For>
        </CardContent>
      </Show>
    </Card>
  );
};

function distancePerGroup(
  sets: (NewTrainingSet | TrainingSet)[]
): DistancePerGroup[] {
  const distances = new Map<GroupEnum, number>();
  let noGroup = 0;

  for (const set of sets) {
    if (!set.group) {
      noGroup += set.totalDistance;
      continue;
    }

    if (!distances.has(set.group)) {
      distances.set(set.group, 0);
    }
    distances.set(
      set.group,
      (distances.get(set.group) || 0) + set.totalDistance
    );
  }

  for (const g of Array.from(distances.keys())) {
    distances.set(g, distances.get(g)! + noGroup);
  }

  return Array.from(distances.entries()).map(([g, d]) => {
    return { group: g, distance: d };
  });
}

export default TrainingSummary;
