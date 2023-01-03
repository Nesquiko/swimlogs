create table if not exists session (
	id uuid primary key,
	created_at timestamp not null,
	modified_at timestamp not null,
	day day not null,
	startTime char(5) not null,
	duration smallint not null
);
