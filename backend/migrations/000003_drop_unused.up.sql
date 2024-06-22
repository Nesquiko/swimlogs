alter table trainings
    drop column if exists total_distance;

alter table trainings
    drop constraint if exists trainings_total_distance_check;

alter table sets
    drop column if exists total_distance;

alter table sets
    add constraint sets_set_order_check check ( set_order >= 0 );

alter table sets
    alter start_type drop not null;

update sets
set start_type = null
where start_type = 'None';

create type new_start_type as enum ('Interval', 'Pause');

alter table sets
    drop constraint if exists sets_rule_check;

alter table sets
    alter column start_type type new_start_type using start_type::text::new_start_type;

drop type start_type;

alter type new_start_type rename to start_type;

alter table sets
    add constraint sets_rule_check check (start_type is null or (start_seconds is not null and start_seconds > 0));
