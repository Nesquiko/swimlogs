import { Component, createSignal, For, Match, Switch } from 'solid-js';
import { createStore, SetStoreFunction } from 'solid-js/store';
import { NewSetComponent, NewTrainingSet, Style } from '~/api/generated';
import { useAppState } from '~/AppContext';
import OptionsGrid from './OptionsGrid';
import { DistanceRange, DISTANCES } from '../common';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '~/components/ui/tabs';
import { Button } from '~/components/ui/button';
import BottomIcons from './BottomIcons';
import {
  IconCheck,
  IconDraggable,
  IconMenu,
  IconPen,
  IconPlus,
  IconTrash,
  IconX,
} from '~/components/icons';
import { Card } from '~/components/ui/card';
import {
  closestCenter,
  createSortable,
  DragDropProvider,
  DragDropSensors,
  DragEvent,
  DragOverlay,
  SortableProvider,
  useDragDropContext,
} from '@thisbeyond/solid-dnd';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '~/components/ui/dropdown-menu';

interface CompoundSetProps {
  set: NewTrainingSet;
  setSet: SetStoreFunction<NewTrainingSet>;
  onSubmitSet: () => void;
  onCancel: () => void;
}

export const CompoundSetForm: Component<CompoundSetProps> = (props) => {
  const [screen, setScreen] = createSignal<'new-component' | 'preview'>(
    (props.set.components?.length ?? 0) === 0 ? 'new-component' : 'preview'
  );

  let [component, setComponent] = createStore(defaultNewComponent(props.set));

  const onAddComponent = () => {
    props.setSet('components', (comps) => [...(comps ?? []), component]);
    [component, setComponent] = createStore(defaultNewComponent(props.set));
    setScreen('preview');
  };

  return (
    <div>
      <Switch>
        <Match when={screen() === 'preview'}>
          <div class="flex flex-col items-center gap-4">
            <CompoundSetPreview
              components={props.set.components ?? []}
              setSet={props.setSet}
            />
            <Button
              size="icon"
              class="rounded-full"
              onClick={() => setScreen('new-component')}
            >
              <IconPlus class="size-8" />
            </Button>
          </div>
        </Match>
        <Match when={screen() === 'new-component'}>
          <ComponentForm component={component} setComponent={setComponent} />
        </Match>
      </Switch>

      <BottomIcons
        left={{
          hidden: screen() !== 'new-component',
          comp: IconX,
          onClick: props.onCancel,
        }}
        right={{
          comp: screen() === 'new-component' ? IconPlus : IconCheck,
          onClick: () => {
            if (screen() === 'new-component') {
              onAddComponent();
            } else {
              props.onSubmitSet();
            }
          },
        }}
      />
    </div>
  );
};

interface CompoundSetPreviewProps {
  components: Array<NewSetComponent>;
  setSet: SetStoreFunction<NewTrainingSet>;
}

const CompoundSetPreview: Component<CompoundSetPreviewProps> = (props) => {
  const ids = () => props.components.map(hashNewSetComponent);

  function throwIfWrongId(id: string | number | null): asserts id is string {
    if (typeof id === 'number') {
      throw new Error('invalid id, it was number and not string');
    }
  }

  const getDraggedComp = (draggedId: string) => {
    throwIfWrongId(draggedId);
    return props.components.find((c) => hashNewSetComponent(c) === draggedId)!;
  };

  const onDragEnd = ({ draggable, droppable }: DragEvent) => {
    if (!draggable || !droppable) {
      return;
    }

    throwIfWrongId(draggable.id);
    throwIfWrongId(droppable.id);

    const fromIndex = ids().indexOf(draggable.id);
    const toIndex = ids().indexOf(droppable.id);

    if (fromIndex === toIndex) {
      return;
    }

    const updatedItems = props.components.slice();
    updatedItems.splice(toIndex, 0, ...updatedItems.splice(fromIndex, 1));
    props.setSet('components', updatedItems);
  };

  const style = (component: NewSetComponent) => {
    let style = Styles.find((s) => s.id === component.styleId);
    if (style) return style;
    style = PullKick.find((s) => s.id === component.styleId);
    if (style) return style;
    style = Exercises.find((s) => s.id === component.styleId);
    if (style) return style;

    throw new Error(`unknown style id ${component.styleId}`);
  };

  // TODO here, implement onComponentDelete so that all the componentOrders
  // are recalculated, also do it on move, then complete onComponentEdit
  const onComponentDelete = (idx: number) => {
    props.setSet('components', (components) =>
      components?.filter((_, i) => i !== idx)
    );
  };

  const onComponentEdit = (idx: number) => {};

  return (
    <DragDropProvider onDragEnd={onDragEnd} collisionDetector={closestCenter}>
      <DragDropSensors />
      <div class="grid auto-cols-fr grid-flow-row auto-rows-min gap-4 self-stretch">
        <SortableProvider ids={ids()}>
          <For each={props.components}>
            {(comp, idx) => (
              <DraggableComponentCard
                id={ids()[idx()]}
                component={comp}
                style={style(comp)}
                onDelete={() => onComponentDelete(idx())}
                onEdit={() => onComponentEdit(idx())}
              />
            )}
          </For>
        </SortableProvider>
      </div>
      <DragOverlay>
        {(active) => {
          if (!active) return;
          throwIfWrongId(active.id);
          const draggedComp = getDraggedComp(active.id);
          return (
            <ComponentCard component={draggedComp} style={style(draggedComp)} />
          );
        }}
      </DragOverlay>
    </DragDropProvider>
  );
};

interface ComponentFormProps {
  component: NewSetComponent;
  setComponent: SetStoreFunction<NewSetComponent>;
}

const ComponentForm: Component<ComponentFormProps> = (props) => {
  const { t } = useAppState();

  const isDistanceValid = () =>
    props.component.distanceMeters >= DistanceRange.start &&
    props.component.distanceMeters <= DistanceRange.end;

  const stylesTabsContent = (tab: string, styles: Array<Style>) => (
    <TabsContent value={tab}>
      <div class="grid grid-cols-3 gap-2">
        <For each={styles}>
          {(style) => (
            <Button
              variant={
                props.component.styleId === style.id ? 'outline' : 'ghost'
              }
              onClick={() => props.setComponent('styleId', style.id)}
            >
              {style.execriseName ?? style.name}
            </Button>
          )}
        </For>
      </div>
    </TabsContent>
  );

  return (
    <div class="flex flex-col gap-8">
      <span class="text-xl">
        <span>{t('newtraining.new.set.part')}</span>
      </span>

      <OptionsGrid
        label={t('general.training.distance')}
        firstRow={[DISTANCES[0], DISTANCES[1], DISTANCES[2]]}
        secondRow={[DISTANCES[3], DISTANCES[4], DISTANCES[5]]}
        value={props.component.distanceMeters}
        onChange={(dist) => props.setComponent('distanceMeters', dist)}
        orDifferentInput={{
          isValid: isDistanceValid(),
          minValue: DistanceRange.start,
          maxValue: DistanceRange.end,
          step: 25,
        }}
      />

      <Tabs defaultValue="style" class="w-full">
        <TabsList class="grid w-full grid-cols-3">
          <TabsTrigger value="style">{t('newtraining.style')}</TabsTrigger>
          <TabsTrigger value="pull-kick">
            {t('newtraining.pull-kick')}
          </TabsTrigger>
          <TabsTrigger value="exercise">
            {t('newtraining.exercise')}
          </TabsTrigger>
        </TabsList>
        {stylesTabsContent('style', Styles)}
        {stylesTabsContent('pull-kick', PullKick)}
        {stylesTabsContent('exercise', Exercises)}
      </Tabs>
    </div>
  );
};

interface DraggableComponentCardProps {
  id: string;
  component: NewSetComponent;
  style: Style;

  onDelete: () => void;
  onEdit: () => void;
}

const DraggableComponentCard: Component<DraggableComponentCardProps> = (
  props
) => {
  const sortable = createSortable(props.id);
  const [state] = useDragDropContext()!;

  return (
    <div
      ref={sortable.ref}
      classList={{
        'opacity-25': sortable.isActiveDraggable,
        'transition-transform': !!state.active.draggable,
      }}
    >
      <ComponentCard
        component={props.component}
        style={props.style}
        asMenuItem={{
          sortable,
          onDelete: props.onDelete,
          onEdit: props.onEdit,
        }}
      />
    </div>
  );
};

interface ComponentCardProps {
  component: NewSetComponent;
  style: Style;

  asMenuItem?: {
    sortable: ReturnType<typeof createSortable>;
    onDelete: () => void;
    onEdit: () => void;
  };
}

const ComponentCard: Component<ComponentCardProps> = (props) => {
  const { t } = useAppState();

  return (
    <Card class="w-full">
      <div class="grid grid-cols-4 items-center justify-between p-4">
        <p class="text-left">
          <IconDraggable
            {...props.asMenuItem?.sortable.dragActivators}
            class="inline size-6 cursor-move touch-none"
          />
        </p>
        <span class="text-lg">{props.component.distanceMeters}m</span>
        <span class="text-lg">
          {props.style.execriseName ?? props.style.name}
        </span>
        <p class="text-right">
          <DropdownMenu>
            <DropdownMenuTrigger disabled={!props.asMenuItem}>
              <IconMenu class="inline size-6 cursor-pointer" />
            </DropdownMenuTrigger>
            <DropdownMenuContent>
              <DropdownMenuItem onClick={props.asMenuItem?.onEdit}>
                <IconPen />
                {t('general.edit')}
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem
                class="text-destructive"
                onClick={props.asMenuItem?.onDelete}
              >
                <IconTrash />
                {t('general.delete')}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </p>
      </div>
    </Card>
  );
};

const defaultNewComponent = (set: NewTrainingSet): NewSetComponent => {
  return {
    repeat: 1,
    distanceMeters: 50,
    componentOrders: [set.components?.length ?? 0],
  };
};

const Styles: Array<Style> = [
  { id: '1', name: 'choice' },
  { id: '2', name: 'medley' },
  { id: '3', name: 'freestyle' },
  { id: '4', name: 'breastroke' },
  { id: '5', name: 'butterfly' },
  { id: '6', name: 'backstroke' },
];

const PullKick: Array<Style> = [
  { id: 'a', name: 'medley', execriseName: 'medley kick' },
  { id: 'b', name: 'medley', execriseName: 'medley pull' },
  { id: 'c', name: 'freestyle', execriseName: 'free kick' },
  { id: 'd', name: 'freestyle', execriseName: 'free pull' },
  { id: 'e', name: 'breastroke', execriseName: 'breastroke kick' },
  { id: 'f', name: 'breastroke', execriseName: 'breastroke pull' },
  { id: 'g', name: 'butterfly', execriseName: 'dolphin kick' },
  { id: 'h', name: 'butterfly', execriseName: 'butterfly pull' },
  { id: 'i', name: 'backstroke', execriseName: 'backstroke kick' },
  { id: 'j', name: 'backstroke', execriseName: 'backstroke pull' },
];

const Exercises: Array<Style> = [
  { id: 'A', name: 'freestyle', execriseName: 'catch up' },
  { id: 'B', name: 'freestyle', execriseName: 'fists' },
  { id: 'C', name: 'breastroke', execriseName: 'on back' },
  { id: 'D', name: 'butterfly', execriseName: 'side kick' },
  { id: 'E', name: 'backstroke', execriseName: 'one hand' },
];

function hashNewSetComponent(component: NewSetComponent): string {
  const essentialValues = [
    component.repeat,
    component.distanceMeters,
    component.description || '',
    component.start?.type || '',
    (component.equipment || []).sort().join(','),
    component.group || '',
    component.intensity || '',
    component.progression || '',
    component.componentOrders.join(','),
    component.iterationOrder || 0,
    component.styleId || '',
  ];

  const stringToHash = essentialValues.join('|');

  let hash = 0;
  for (let i = 0; i < stringToHash.length; i++) {
    const char = stringToHash.charCodeAt(i);
    hash = (hash << 5) - hash + char;
    hash = hash & hash;
  }

  return Math.abs(hash).toString(16);
}
