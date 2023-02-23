create table if not exists session (
	id uuid primary key,
	created_at timestamp not null,
	modified_at timestamp not null,
	version integer not null,

	day day not null,
	startTime time not null,
	duration smallint not null,

	constraint session_duration_check check (duration > 0)
);

