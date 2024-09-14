import { Component, For } from 'solid-js';
import { Button } from './ui/button';
import { useAppState } from '~/AppContext';
import { Dictionary } from './i18n';
import { IconCalendar, IconHome, IconUser } from './icons';
import { IconProps } from 'solid-icons';
import { useLocation, useNavigate } from '@solidjs/router';

interface NavItem {
  labelKey: Extract<keyof Dictionary, `nav.${any}`>;
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

  const navItem = (item: NavItem, i: number) => {
    return (
      <Button
        variant="outline"
        class="flex w-full flex-col rounded-none border-0 border-t py-6"
        classList={{
          'text-primary hover:text-primary':
            location.pathname === item.location,
          'hover:text-accent-foreground': location.pathname !== item.location,
          'rounded-tl-2xl': i === 0,
          'rounded-tr-2xl': i === NavItems.length - 1,
        }}
        onClick={() => navigate(item.location, { replace: true })}
      >
        <item.icon />
        <span>{t(item.labelKey)}</span>
      </Button>
    );
  };

  return (
    <nav class="fixed bottom-0 left-0 right-0 flex w-full justify-center bg-background">
      <For each={NavItems}>{(item, i) => navItem(item, i())}</For>
    </nav>
  );
};

export default BottomNav;
