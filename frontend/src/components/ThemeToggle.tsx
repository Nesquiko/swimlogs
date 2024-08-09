import { Button } from '~/components/ui/button';
import { IconMoon, IconSun } from './icons';
import { useAppState } from '~/AppContext';

export function ThemeToggle() {
  const context = useAppState();

  return (
    <Button variant="outline" onClick={() => context.setDark(!context.isDark)}>
      {context.isDark ? (
        <IconMoon class="size-6 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
      ) : (
        <IconSun class="size-6 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
      )}
    </Button>
  );
}

export default ThemeToggle;
