create table if not exists set(
	id uuid primary key,

	num smallint not null,
	repeat smallint not null,
	distance smallint not null,
	what text not null,
	starting_rule starting_rule not null,
	rule_seconds smallint,
	total_dist smallint not null,
	block_id uuid references block on delete cascade,

	constraint set_repeat_check check (repeat > 0),
	constraint set_distance_check check (distance > 0),
	constraint set_rule_check check (not (rule_seconds is not null) or rule_seconds > 0), -- equivalent with implication: if x then y
	constraint set_total_dist_check check (total_dist > 0)
);
