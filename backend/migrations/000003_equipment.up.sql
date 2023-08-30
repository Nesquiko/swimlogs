create type equipment as enum ('Fins', 'Monofin', 'Snorkel', 'Board', 'Paddles');

alter table sets
    add column if not exists equipment equipment[];