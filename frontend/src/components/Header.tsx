import { A, useLocation } from '@solidjs/router';
import {
  Component,
  createEffect,
  createSignal,
  For,
  Match,
  Show,
  Switch,
} from 'solid-js';
import { useOnBackcontext } from '../pages/Routing';
import { openDrawer } from './Drawer';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from './ui/dropdown-menu';

interface MenuItem {
  text: string;
  icon?: string;
  textColorCls?: string;
  onClick: () => void;
}

const [headerMenu, setHeaderMenu] = createSignal<
  { items: MenuItem[]; lastItem?: MenuItem } | undefined
>(undefined);

const clearHeaderMenu = () => setHeaderMenu(undefined);

const Header: Component = () => {
  const [headerState, setHeaderState] = createSignal<{
    state: 'menu' | 'back';
    onBack?: () => void;
  }>({
    state: 'menu',
  });
  const location = useLocation();
  const [onBack] = useOnBackcontext();

  createEffect(() => {
    if (location.pathname === '/') {
      setHeaderState({ state: 'menu' });
    } else {
      setHeaderState({
        state: 'back',
        onBack,
      });
    }
  });

  const menuItem = (item: MenuItem) => {
    const textColorCls = item.textColorCls ?? 'text-gray-700';
    return (
      <DropdownMenuItem onClick={item.onClick}>
        {item.icon && (
          <i
            classList={{ [item.icon]: true, [textColorCls]: true }}
            class="fa-solid fa-pen pr-2"
          />
        )}
        <span>{item.text}</span>
      </DropdownMenuItem>
    );
  };

  return (
    <div class="w-full bg-sky-500 px-4 py-2 grid grid-cols-5 items-center">
      <Switch>
        <Match when={headerState().state === 'menu'}>
          <i
            class="col-start-1 fa-solid fa-bars fa-2xl cursor-pointer text-white"
            onClick={openDrawer}
          ></i>
        </Match>
        <Match when={headerState().state === 'back'}>
          <i
            class="col-start-1 fa-solid fa-arrow-left fa-2xl cursor-pointer text-white"
            onClick={headerState().onBack}
          ></i>
        </Match>
      </Switch>
      <A
        href="/"
        class="col-span-3 col-start-2 text-center cursor-pointer text-xl font-bold text-white"
      >
        SwimLogs
      </A>

      <Show when={headerMenu()} keyed>
        {(menu) => {
          return (
            <DropdownMenu>
              <DropdownMenuTrigger>
                <i class="inline-block w-full text-right fa-solid col-start-5 fa-ellipsis fa-2xl cursor-pointer text-white" />
              </DropdownMenuTrigger>
              <DropdownMenuContent>
                <For each={menu.items}>{(item) => menuItem(item)}</For>
                {menu.lastItem && (
                  <>
                    <DropdownMenuSeparator />
                    {menuItem(menu.lastItem)}
                  </>
                )}
              </DropdownMenuContent>
            </DropdownMenu>
          );
        }}
      </Show>
    </div>
  );
};

export default Header;
export { setHeaderMenu, clearHeaderMenu };
