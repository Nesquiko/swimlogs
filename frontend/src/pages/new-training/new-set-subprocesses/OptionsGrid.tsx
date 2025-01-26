import { Component, For, Show } from 'solid-js';
import { useAppState } from '~/AppContext';
import { Button } from '~/components/ui/button';
import { Label } from '~/components/ui/label';
import { MobileNumberField } from '~/components/ui/number-field';
import { randomId } from '~/lib/str';

interface OptionsGridProps {
  label: string;
  firstRow: readonly [number, number, number];
  secondRow: readonly [number, number, number];

  value: number;
  onChange: (val: number) => void;

  orDifferentInput: {
    isValid: boolean;
    minValue?: number;
    maxValue?: number;
    step?: number;
  };
}

const OptionsGrid: Component<OptionsGridProps> = (props: OptionsGridProps) => {
  const { t } = useAppState();
  const id = randomId();

  const item = (val: number) => {
    return (
      <Button
        variant={props.value === val ? 'outline' : 'ghost'}
        onClick={() => props.onChange(val)}
      >
        {val}
      </Button>
    );
  };

  return (
    <div class="grid grid-cols-3">
      <Label class="col-span-3 pb-4" for={`${id}-input`}>
        {props.label}
      </Label>
      <For each={props.firstRow}>{item}</For>
      <For each={props.secondRow}>{item}</For>
      <div class="relative col-span-3 flex items-center py-1 text-sm">
        <div class="flex-grow border-t"></div>
        <span class="mx-2 text-muted-foreground">
          {t('newtraining.or.different')}
        </span>
        <div class="flex-grow border-t"></div>
      </div>
      <MobileNumberField
        id={id}
        class="col-span-3 w-full"
        value={props.value}
        minValue={props.orDifferentInput.minValue}
        maxValue={props.orDifferentInput.maxValue}
        step={props.orDifferentInput.step ?? 1}
        onRawValueChange={props.onChange}
        isValid={props.orDifferentInput.isValid}
        changeOnWheel
      />
    </div>
  );
};

export default OptionsGrid;
