import { useNavigate } from '@solidjs/router';
import { For, JSX, ParentComponent, Show, type Component } from 'solid-js';
import { TrainingSummary } from 'swimlogs-api';
import { useAppState } from '~/AppContext';
import {
  IconBug,
  IconClock,
  IconRepeat,
  IconStopwatch,
  IconX,
} from '~/components/icons';
import { Alert, AlertDescription, AlertTitle } from '~/components/ui/alert';
import { Card, CardContent, CardHeader, CardTitle } from '~/components/ui/card';
import { Separator } from '~/components/ui/separator';
import { Skeleton } from '~/components/ui/skeleton';
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
    <TodaySectionCard
      title={
        <>
          {t('home.today.title')}{' '}
          <span class="font-normal">{formatDate(new Date(), true)}</span>
        </>
      }
    >
      <For each={props.trainings}>
        {(training) => (
          <TodaysTraining
            training={training}
            onClick={(id) => navigate('/training/' + id)}
          />
        )}
      </For>
    </TodaySectionCard>
  );
};

export const TodaySectionError: Component = () => {
  const { t } = useAppState();

  return (
    <Alert variant="destructive">
      <IconBug />
      <AlertTitle>{t('messages.failed.today.summaries.title')}</AlertTitle>
      <AlertDescription>
        {t('messages.failed.today.summaries.description')}
      </AlertDescription>
    </Alert>
  );
};

export const TodaySectionLoading: Component = () => {
  return (
    <TodaySectionCard title={<Skeleton height={16} width={120} radius={10} />}>
      <div class="flex w-full flex-col gap-2 px-6 py-2">
        <Separator />
        <span class="flex items-center gap-4">
          <Skeleton height={16} width={100} radius={10} />
          <Skeleton height={16} width={124} radius={10} />
        </span>
      </div>
    </TodaySectionCard>
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
      <div class="flex w-full flex-col gap-2 px-6 py-2">
        <For each={props.training.mainSets}>
          {(mainSet) => (
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
          )}
        </For>
      </div>
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
        {mainSets()}
      </div>
    </div>
  );
};

interface TodaySectionCardProps {
  title?: string | JSX.Element;
}

const TodaySectionCard: ParentComponent<TodaySectionCardProps> = (props) => {
  return (
    <Card class="w-full max-w-lg">
      <Show when={props.title}>
        <CardHeader>
          <CardTitle>{props.title}</CardTitle>
        </CardHeader>
      </Show>
      <Show when={props.children}>
        <CardContent>{props.children}</CardContent>
      </Show>
    </Card>
  );
};

export default TodaySection;
