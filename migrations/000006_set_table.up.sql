create table if not exists set(
	id uuid primary key,
	created_at timestamp not null,
	modified_at timestamp not null,

	repeat smallint not null,
	distance smallint not null,
	what text not null,
	starting_rule starting_rule not null,
	rule_seconds smallint,

	constraint set_repeat_check check (repeat > 0),
	constraint set_distance_check check (distance > 0),
	constraint set_rule_check check (not (rule_seconds is not null) or rule_seconds > 0) -- equivalent with implication: if x then y
);
