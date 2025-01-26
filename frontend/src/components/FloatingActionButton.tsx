import { Component, Show } from 'solid-js';
import { Button } from './ui/button';
import { IconProps } from './icons';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from './ui/dropdown-menu';

interface FloatingActionButtonProps {
  label?: string;
  icon: Component<IconProps>;
}

const FloatingActionButton: Component<FloatingActionButtonProps> = (props) => {
  const size = props.label ? 'default' : 'icon';

  return (
    <DropdownMenu>
      <DropdownMenuTrigger class="fixed bottom-16 right-4">
        <Button size={size}>
          <props.icon size={30} />
          <Show when={props.label}>
            <span class="pl-2">{props.label}</span>
          </Show>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent class="border-0">
        <DropdownMenuItem>
          <Button>
            <span class="pl-2">Normal</span>
          </Button>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Push</span>
        </DropdownMenuItem>
        <DropdownMenuItem>
          <span>Update Project</span>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default FloatingActionButton;
