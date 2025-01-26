import { Component, Show } from 'solid-js';
import { IconProps } from '~/components/icons';

interface Icon {
  comp: Component<IconProps>;
  onClick: () => void;
  hidden?: boolean;
}

interface BottomIconsProps {
  left?: Icon;
  right?: Icon;
}

const BottomIcons: Component<BottomIconsProps> = (props) => {
  return (
    <div class="fixed bottom-4 left-4 right-4 grid grid-cols-2">
      <Show when={props.left} keyed>
        {(left) => (
          <p
            classList={{ invisible: left.hidden }}
            class="h-full text-left text-3xl leading-none"
          >
            <left.comp class="inline cursor-pointer" onClick={left.onClick} />
          </p>
        )}
      </Show>
      <Show when={props.right} keyed>
        {(right) => (
          <p
            classList={{ invisible: right.hidden }}
            class="col-start-2 h-full text-right text-3xl leading-none"
          >
            <right.comp class="inline cursor-pointer" onClick={right.onClick} />
          </p>
        )}
      </Show>
    </div>
  );
};

export default BottomIcons;
