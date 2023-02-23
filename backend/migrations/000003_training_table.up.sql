create table if not exists training(
	id uuid primary key,
	created_at timestamp not null,
	modified_at timestamp not null,
	version integer not null,

	date date not null,
	day day not null,
	startTime time not null,
	duration smallint not null,
	total_dist smallint not null,

	constraint training_duration_check check (duration > 0),
	constraint training_total_dist_check check (total_dist > 0)
);
