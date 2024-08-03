create type new_start_type as enum ('interval', 'pause');

alter table sets
    drop constraint if exists sets_rule_check;

alter table sets
    alter column start_type type new_start_type using start_type::text::new_start_type;

drop type start_type;

alter type new_start_type rename to start_type;

alter table sets
    add constraint sets_rule_check check (start_type is null or (start_seconds is not null and start_seconds > 0));


create type new_equipment as enum ('fins', 'monofin', 'snorkel', 'board', 'paddles');

alter table sets
    alter column equipment type new_equipment[] using equipment::equipment[]::text[]::new_equipment[];

drop type equipment;

alter type new_equipment rename to equipment;
