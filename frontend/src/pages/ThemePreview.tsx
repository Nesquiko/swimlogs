import { type Component } from 'solid-js';
import { ThemeToggleButton } from './profile/AppearanceSection';
import { useAppState } from '~/AppContext';

const ThemePreview: Component = () => {
  const context = useAppState();

  return (
    <div class="flex flex-col gap-4 p-2">
      <ThemeToggleButton context={context} />
      <p class="bg-background text-xl text-foreground">
        Background + Foreground
      </p>
      <p class="bg-muted text-xl text-muted-foreground">Muted</p>
      <p class="bg-popover text-xl text-popover-foreground">Popover</p>
      <p class="bg-border text-xl text-input">Border + Input</p>
      <p class="bg-card text-xl text-card-foreground">Card</p>
      <p class="bg-primary text-xl text-primary-foreground">Primary</p>
      <p class="bg-secondary text-xl text-secondary-foreground">Secondary</p>
      <p class="bg-accent text-xl text-accent-foreground">Accent</p>
      <p class="bg-destructive text-xl text-destructive-foreground">
        Destructive
      </p>
      <p class="bg-info text-xl text-info-foreground">Info</p>
      <p class="bg-success text-xl text-success-foreground">Success</p>
      <p class="bg-warning text-xl text-warning-foreground">Warning</p>
      <p class="bg-error text-xl text-error-foreground">Error</p>
      <p class="bg-ring text-xl text-ring">Ring</p>
    </div>
  );
};

export default ThemePreview;
