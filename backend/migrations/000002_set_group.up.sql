create type set_group as enum ('bifi', 'long', 'middle', 'mono', 'sprint');
alter table sets add column if not exists "group" set_group;

