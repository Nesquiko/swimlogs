import { Component, For } from 'solid-js';
import { NewTrainingSet } from '~/api/generated';
import { MobileNumberField } from '~/components/ui/number-field';
import { SetStoreFunction } from 'solid-js/store';
import { useAppState } from '~/AppContext';
import { randomId } from '~/lib/str';
import { Button } from '~/components/ui/button';
import { Label } from '~/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs';
import AllSameSetsForm from './AllSameSetsForm';
import { Checkbox } from '~/components/ui/checkbox';
import { DistanceRange, DISTANCES, RepeatRange } from './common';

interface NormalSetFormProps {
  set: NewTrainingSet;
  setSet: SetStoreFunction<NewTrainingSet>;
}

const NormalSetForm: Component<NormalSetFormProps> = (props) => {
  const { t } = useAppState();
  const distanceId = randomId();

  const isRepeatValid = () =>
    props.set.repeat >= RepeatRange.start &&
    props.set.repeat <= RepeatRange.end;
  const isDistanceValid = () =>
    props.set.distanceMeters >=
      (DistanceRange.validationStart ?? DistanceRange.start) &&
    props.set.distanceMeters <= DistanceRange.end;

  const distanceItem = (dist: number) => {
    return (
      <Button
        variant={props.set.distanceMeters === dist ? 'outline' : 'ghost'}
        onClick={() => props.setSet('distanceMeters', dist)}
      >
        {dist}
      </Button>
    );
  };

  return (
    <div>
      <div class="flex items-center justify-end gap-2 pb-4">
        <Checkbox
          class="space-x-0"
          checked={props.set.isMain}
          onChange={(v) => props.setSet('isMain', v)}
          id="ismain"
        />
        <Label for="ismain-input">{t('newtraining.is.main')}</Label>
      </div>
      <div class="grid grid-cols-2">
        <MobileNumberField
          value={props.set.repeat}
          minValue={RepeatRange.start}
          maxValue={RepeatRange.end}
          label={t('newtraining.repeat.count')}
          onRawValueChange={(repeat) => props.setSet('repeat', repeat)}
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
            value={props.set.distanceMeters}
            minValue={DistanceRange.start}
            maxValue={DistanceRange.end}
            step={25}
            onRawValueChange={(dist) => props.setSet('distanceMeters', dist)}
            isValid={isDistanceValid()}
          />
        </div>
      </div>

      <Label>{t('general.training.repeat')}</Label>
      <Tabs defaultValue="same">
        <TabsList class="grid w-full grid-cols-2">
          <TabsTrigger value="same">
            {t('newtraining.repetitions.all.same')}
          </TabsTrigger>
          <TabsTrigger disabled={props.set.repeat <= 1} value="different">
            {t('newtraining.repetitions.different')}
          </TabsTrigger>
        </TabsList>
        <div class="pt-2">
          <TabsContent value="same">
            <AllSameSetsForm set={props.set} setSet={props.setSet} />
          </TabsContent>
        </div>
      </Tabs>
    </div>
  );
};

export default NormalSetForm;
