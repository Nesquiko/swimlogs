import { useNavigate } from '@solidjs/router';
import { For, Show, type Component } from 'solid-js';
import { TrainingSummary } from 'swimlogs-api';
import { useAppState } from '~/AppContext';
import {
  IconClock,
  IconRepeat,
  IconStopwatch,
  IconX,
} from '~/components/icons';
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card';
import { Separator } from '~/components/ui/separator';
import {
  formatDate,
  formatTime,
  secondsToMintesAndSeconds,
} from '~/lib/datetime';

interface TodaySectionProps {
  trainings: Array<TrainingSummary>;
}

const TodaySection: Component<TodaySectionProps> = (props) => {
  const navigate = useNavigate();
  const { t } = useAppState();

  return (
    <Card class="w-full max-w-lg">
      <CardHeader>
        <CardTitle>
          {t('home.today.title')}{' '}
          <span class="font-normal">{formatDate(new Date(), true)}</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <For each={props.trainings}>
          {(training) => (
            <TodaysTraining
              training={training}
              onClick={(id) => navigate('/training/' + id)}
            />
          )}
        </For>
      </CardContent>
    </Card>
  );
};

interface TodaysTrainingProps {
  training: TrainingSummary;
  onClick: (id: string) => void;
}

const TodaysTraining: Component<TodaysTrainingProps> = (props) => {
  const { t } = useAppState();

  const mainSets = () => {
    return (
      <For each={props.training.mainSets}>
        {(mainSet) => (
          <p>
            <span class="flex items-center gap-4">
              <span class="flex items-center gap-1">
                <IconRepeat class="inline text-primary" />
                {mainSet.repeat}
                <IconX class="inline text-primary" />
                {mainSet.distanceMeters}m
              </span>
              <Show when={mainSet.startType} keyed>
                {(startType) => (
                  <span class="flex items-center gap-1">
                    <IconStopwatch class="inline text-primary" />
                    <span>{t(`general.starttype.${startType}`)}</span>
                    <span>
                      {secondsToMintesAndSeconds(mainSet.startSeconds)}
                    </span>
                  </span>
                )}
              </Show>
            </span>
          </p>
        )}
      </For>
    );
  };

  return (
    <div>
      <Separator />
      <div
        class="cursor-pointer py-4"
        onClick={() => props.onClick(props.training.id)}
      >
        <div class="flex items-center justify-between text-lg">
          <span class="flex items-center gap-2">
            <IconClock class="inline text-primary" />
            {formatTime(props.training.start)}
          </span>
          <span>{props.training.totalDistance / 1000} km</span>
        </div>
        <div class="flex w-full flex-col gap-2 px-6 py-2">{mainSets()}</div>
      </div>
    </div>
  );
};

export default TodaySection;
