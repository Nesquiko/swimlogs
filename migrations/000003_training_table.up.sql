create table if not exists training(
	id uuid primary key,
	created_at timestamp not null,
	modified_at timestamp not null,

	date date not null,
	day day not null,
	startTime time not null,
	duration smallint not null,

	constraint training_duration_check check (duration > 0)
);
