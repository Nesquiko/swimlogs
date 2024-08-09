import { Component, For } from 'solid-js';
import { Button } from './ui/button';
import { useAppState } from '~/AppContext';
import { Dictionary } from './i18n';
import { IconCalendar, IconHome, IconUser } from './icons';
import { IconProps } from 'solid-icons';
import { useLocation, useNavigate } from '@solidjs/router';

interface NavItem {
  labelKey: keyof Dictionary;
  icon: Component<IconProps>;
  location: string;
}

const NavItems: Array<NavItem> = [
  { labelKey: 'nav.home', icon: IconHome, location: '/' },
  {
    labelKey: 'nav.calendar',
    icon: IconCalendar,
    location: '/calendar',
  },
  { labelKey: 'nav.profile', icon: IconUser, location: '/profile' },
];

const BottomNav: Component = () => {
  const { t } = useAppState();
  const location = useLocation();
  const navigate = useNavigate();

  const navItem = (item: NavItem) => {
    return (
      <Button
        variant="outline"
        class="flex w-full flex-col rounded-none py-6"
        classList={{
          'text-primary hover:text-primary':
            location.pathname === item.location,
          'hover:text-accent-foreground': location.pathname !== item.location,
        }}
        onClick={() => navigate(item.location)}
      >
        <item.icon />
        <span>{t(item.labelKey)}</span>
      </Button>
    );
  };

  return (
    <nav class="fixed bottom-0 flex w-full justify-center">
      <For each={NavItems}>{navItem}</For>
    </nav>
  );
};

export default BottomNav;
