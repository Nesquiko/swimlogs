create table if not exists set_components
(
    id               uuid primary key default gen_random_uuid(),
    set_id           uuid       not null references sets on delete cascade,

    component_orders smallint[] not null,
    iteration_order  smallint,
    repeat           smallint   not null,
    distance_meters  smallint   not null,

    start_type       start_type,
    start_seconds    smallint,
    style_id         uuid       references styles on delete set null,
    intensity        set_intensity,
    progression      set_progression,
    description      text,
    equipment        equipment[],
    "group"          set_group,

    constraint set_components_repeat_check check (repeat > 0),
    constraint set_components_components_orders_check check (array_length(component_orders, 1) > 0),
    constraint set_components_distance_check check (distance_meters is null or distance_meters > 0)
);
