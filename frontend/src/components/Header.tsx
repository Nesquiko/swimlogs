import { A, useLocation } from '@solidjs/router';
import { Component, createEffect, createSignal, Match, Switch } from 'solid-js';
import { useOnBackcontext } from '../pages/Routing';
import { openDrawer } from './Drawer';

type HeaderState = {
  state: 'menu' | 'back';
  onBack?: () => void;
};

const Header: Component = () => {
  const [headerState, setHeaderState] = createSignal<HeaderState>({
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

  return (
    <div id="topbar" class="w-full bg-sky-500 p-2">
      <div class="inline-block w-1/5 align-middle">
        <Switch>
          <Match when={headerState().state === 'menu'}>
            <i
              class="fa-solid fa-bars fa-2xl cursor-pointer text-white"
              onClick={openDrawer}
            ></i>
          </Match>
          <Match when={headerState().state === 'back'}>
            <i
              class="fa-solid fa-arrow-left fa-2xl cursor-pointer text-white"
              onClick={headerState().onBack}
            ></i>
          </Match>
        </Switch>
      </div>
      <div class="inline-block w-3/5 text-center align-middle">
        <A href="/" class="cursor-pointer text-xl font-bold text-white">
          SwimLogs
        </A>
      </div>
    </div>
  );
};

export default Header;
