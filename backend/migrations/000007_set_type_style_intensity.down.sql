alter table sets
    drop column if exists type;

alter table sets
    drop column if exists style_id;

alter table sets
    drop column if exists intensity;

alter table sets
    drop column if exists progression;

drop type if exists set_type;
drop type if exists set_intensity;
drop type if exists set_progression;
