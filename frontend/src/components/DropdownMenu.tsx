import { Dropdown } from 'flowbite';
import { Component, For, onMount, Show } from 'solid-js';
import { randomId } from '../lib/str';

interface MenuItem {
  text: string;
  icon: string;
  onClick: () => void;
  disabled?: boolean;
}

interface DropdownMenuProps {
  icon: string;
  items: MenuItem[];
  finalItem?: MenuItem;
}

const DropdownMenu: Component<DropdownMenuProps> = (props) => {
  const id = randomId();

  let target: HTMLDivElement;
  let triggerEl: HTMLButtonElement;
  let dropdown: Dropdown;

  onMount(() => {
    dropdown = new Dropdown(target, triggerEl, {}, { id });
  });

  const menuItem = (opts: MenuItem) => {
    if (opts.disabled) {
      return;
    }
    return (
      <a
        class="block px-4 py-2 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white"
        onClick={() => {
          opts.onClick();
          dropdown.hide();
        }}
      >
        <div>
          <i class={`fa-solid ${opts.icon} pr-2`}></i>
          {opts.text}
        </div>
      </a>
    );
  };

  return (
    <div class="inline-block text-left">
      <button
        ref={triggerEl!}
        class="inline-flex items-center text-center text-sm font-medium text-sky-900"
        type="button"
      >
        <i class={`fa-solid ${props.icon} fa-2xl cursor-pointer`}></i>
      </button>
      <div
        id={id}
        ref={target!}
        class="z-10 hidden w-44 divide-y divide-gray-100 rounded-lg bg-white shadow dark:divide-gray-600 dark:bg-gray-700"
      >
        <ul
          class="py-2 text-sm text-gray-700 dark:text-gray-200"
          aria-labelledby="dropdownMenuIconHorizontalButton"
        >
          <For each={props.items}>{(item) => <li>{menuItem(item)}</li>}</For>
        </ul>
        <Show when={props.finalItem}>
          <div class="py-2 text-red-500">{menuItem(props.finalItem!)}</div>
        </Show>
      </div>
    </div>
  );
};

export default DropdownMenu;
