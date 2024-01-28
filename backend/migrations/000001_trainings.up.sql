create table if not exists trainings
(
    id             uuid primary key default gen_random_uuid(),

    start          timestamp with time zone not null,
    duration_min   smallint                 not null,
    total_distance smallint                 not null,

    created_at     timestamp with time zone not null,
    modified_at    timestamp with time zone not null,

    constraint trainings_duration_check check (duration_min > 0),
    constraint trainings_total_distance_check check (total_distance > 0)
);

create type start_type as enum ('None', 'Interval', 'Pause');
create type equipment as enum ('Fins', 'Monofin', 'Snorkel', 'Board', 'Paddles');

create table if not exists sets
(
    id              uuid primary key default gen_random_uuid(),
    training_id     uuid references trainings on delete cascade not null,

    set_order       smallint                                    not null,
    repeat          smallint                                    not null,
    distance_meters smallint                                    not null,
    description     text,
    start_type      start_type                                  not null,
    start_seconds   smallint,
    total_distance  smallint                                    not null,
    equipment       equipment[]

    constraint sets_repeat_check check (repeat > 0),
    constraint sets_distance_check check (distance_meters > 0),
    constraint sets_rule_check check (start_type = 'None' or (start_seconds is not null and start_seconds > 0))
);
