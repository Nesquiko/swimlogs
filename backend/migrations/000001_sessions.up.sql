create type day as enum ('Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday');

create table if not exists sessions
(
    id          uuid primary key default gen_random_uuid(),

    day         day                      not null,
    start_time  time                     not null,
    duration    smallint                 not null,

    created_at  timestamp with time zone not null,
    modified_at timestamp with time zone not null,

    constraint sessions_duration_check check (duration > 0),
    constraint sessions_unique_check unique (day, start_time, duration)
);

