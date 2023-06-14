create table if not exists trainings
(
    id             uuid primary key default gen_random_uuid(),

    date           date      not null,
    start_time     time      not null,
    duration       smallint  not null,
    total_distance smallint  not null,

    created_at     timestamp not null,
    modified_at    timestamp not null,

    constraint trainings_duration_check check (duration > 0),
    constraint trainings_total_distance_check check (total_distance > 0)
);

create table if not exists blocks
(
    id             uuid primary key default gen_random_uuid(),

    num            smallint     not null,
    repeat         smallint     not null,
    name           varchar(255) not null,
    total_distance smallint     not null,
    training_id    uuid references trainings on delete cascade,

    constraint blocks_repeat_check check (repeat > 0),
    constraint blocks_total_distance_check check (total_distance > 0)
);

create type starting_rule as enum ('None', 'Interval', 'Pause');

create table if not exists sets
(
    id             uuid primary key default gen_random_uuid(),

    num            smallint      not null,
    repeat         smallint      not null,
    distance       smallint      not null,
    what           text          not null,
    starting_rule  starting_rule not null,
    rule_seconds   smallint,
    total_distance smallint      not null,
    block_id       uuid references blocks on delete cascade,

    constraint sets_repeat_check check (repeat > 0),
    constraint sets_distance_check check (distance > 0),
    constraint sets_rule_check check (starting_rule = 'None' or (rule_seconds is not null and rule_seconds > 0)),
    constraint sets_total_distance_check check (total_distance = repeat * distance)
);
