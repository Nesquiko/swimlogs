import { Trans, useTransContext } from '@mbarzda/solid-i18next';
import { Component, For, Show } from 'solid-js';
import {
  GroupEnum,
  NewTraining,
  NewTrainingSet,
  Training,
  TrainingSet,
} from 'swimlogs-api';
import { GroupColors } from '../lib/set-groups';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';

interface TrainingSummaryProps {
  training: NewTraining | Training;
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
    <Card>
      <CardHeader
        classList={{ 'rounded-lg': uniqueGroups === 0 }}
        class="p-2 bg-sky-200 rounded-t-lg"
      >
        <CardTitle class="text-xl text-sky-900">
          <span
            classList={{ 'w-3/4': uniqueGroups === 0 }}
            class="inline-block"
          >
            <Trans key={labelKey} />
          </span>
          {uniqueGroups === 0 && (
            <span class="inline-block">{` ${(props.training.totalDistance / 1000).toPrecision(2)}km`}</span>
          )}
        </CardTitle>
      </CardHeader>
      <Show when={uniqueGroups !== 0}>
        <CardContent class="space-y-2 pt-2">
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
