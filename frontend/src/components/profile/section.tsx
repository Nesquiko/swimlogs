import { ParentComponent, Show, type Component } from 'solid-js';
import { IconProps } from '../icons';
import { Separator } from '../ui/separator';

interface SectionHeaderProps {
  header: string;
}

export const SectionHeader: Component<SectionHeaderProps> = (props) => (
  <p class="pb-4 text-xl tracking-tight">{props.header}</p>
);

interface SectionItemProps {
  icon: Component<IconProps>;
  label: string;
  afterLabel?: string;
}

export const SectionItem: ParentComponent<SectionItemProps> = (props) => {
  return (
    <div>
      <div class="flex items-center justify-between py-2">
        <div class="flex items-center gap-2">
          <props.icon class="h-6 w-6 text-primary" />
          <p class="text-lg">
            {props.label}
            <Show when={props.afterLabel}>
              <span class="text-muted-foreground">{props.afterLabel}</span>
            </Show>
          </p>
        </div>
        {props.children}
      </div>
      <Separator />
    </div>
  );
};
