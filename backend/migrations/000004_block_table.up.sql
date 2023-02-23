create table if not exists block(
	id uuid primary key,

	num smallint not null,
	repeat smallint not null,
	name varchar(255) not null,
	training_id uuid references training on delete cascade,
	total_dist smallint not null,

	constraint block_repeat_check check (repeat > 0),
	constraint block_total_dist_check check (total_dist > 0)
);
