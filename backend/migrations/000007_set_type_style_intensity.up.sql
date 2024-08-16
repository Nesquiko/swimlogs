create type set_type as enum ('normal', 'compound', 'pyramid', 'super');
create type set_intensity as enum ('rec', 'en1', 'en2', 'en3', 'sp1', 'sp2', 'sp3');
create type set_progression as enum ('asc', 'desc');

alter table sets
    add column if not exists type set_type not null default 'normal';

alter table sets
    add column if not exists style_id uuid references styles on delete set null;

alter table sets
    add column if not exists intensity set_intensity;

alter table sets
    add column if not exists progression set_progression;


create type new_start_type as enum ('interval', 'pause', 'last-finishes');

alter table sets
    drop constraint if exists sets_rule_check;

alter table sets
    alter column start_type type text using start_type::text;

drop type start_type;

update sets set start_type = lower(start_type);

alter table sets
    alter column start_type type new_start_type using start_type::text::new_start_type;

alter type new_start_type rename to start_type;

