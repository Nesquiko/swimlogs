create type new_start_type as enum ('interval', 'pause');

alter table sets
    drop constraint if exists sets_rule_check;

alter table sets
    alter column start_type type text using start_type::text;

drop type start_type;

update sets
set start_type = lower(start_type);

alter table sets
    alter column start_type type new_start_type using start_type::text::new_start_type;

alter type new_start_type rename to start_type;

create type new_equipment as enum ('fins', 'monofin', 'snorkel', 'board', 'paddles', 'parachute', 'pull buoy');

alter table sets
    alter column equipment type text[] using equipment::equipment[]::text[];

drop type equipment;

update sets
set equipment = lower(equipment::text)::text[];

alter table sets
    alter column equipment type new_equipment[] using equipment::text[]::new_equipment[];

alter type new_equipment rename to equipment;

create type new_set_group as enum ('freestyle', 'breaststroke', 'backstroke', 'butterfly', 'bifi', 'long', 'middle', 'mono', 'sprint');

alter table sets
    alter column "group" type text using "group"::set_group::text;

drop type set_group;

alter table sets
    alter column "group" type new_set_group using "group"::text::new_set_group;

alter type new_set_group rename to set_group;

