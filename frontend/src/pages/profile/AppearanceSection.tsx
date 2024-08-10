import { type Component } from 'solid-js';
import { useAppState } from '~/AppContext';
import {
  IconColorTheme,
  IconLanguage,
  IconMoon,
  IconSun,
} from '~/components/icons';
import { SectionHeader, SectionItem } from '~/components/profile/section';
import { Button } from '~/components/ui/button';
import { capitalize } from '~/lib/str';

const AppearanceSection: Component = () => {
  const context = useAppState();
  const { t } = context;

  return (
    <div>
      <SectionHeader header={t('profile.preferences.headline')} />
      <SectionItem
        label={t('profile.preferences.language.headline')}
        afterLabel={` - ${capitalize(context.locale)}`}
        icon={IconLanguage}
      >
        <div class="space-x-4">
          <Button
            variant={context.locale === 'en' ? 'default' : 'ghost'}
            onClick={() => context.setLocale('en')}
          >
            En
          </Button>
          <Button
            variant={context.locale === 'sk' ? 'default' : 'ghost'}
            onClick={() => context.setLocale('sk')}
          >
            Sk
          </Button>
        </div>
      </SectionItem>

      <SectionItem
        icon={IconColorTheme}
        label={t('profile.preferences.theme.headline')}
        afterLabel={` - ${context.isDark ? t('profile.preferences.theme.dark') : t('profile.preferences.theme.light')}`}
      >
        <Button
          aria-label={t('profile.preferences.theme.headline')}
          variant="ghost"
          onClick={() => context.setDark(!context.isDark)}
        >
          {context.isDark ? (
            <IconMoon class="size-6 rotate-90 scale-0 transition-all dark:rotate-0 dark:scale-100" />
          ) : (
            <IconSun class="size-6 rotate-0 scale-100 transition-all dark:-rotate-90 dark:scale-0" />
          )}
        </Button>
      </SectionItem>
    </div>
  );
};

export default AppearanceSection;
