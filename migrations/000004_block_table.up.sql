create table if not exists block(
	id uuid primary key,

	repeat smallint not null,
	name varchar(255) not null,
	training_id uuid references training on delete cascade,

	constraint block_repeat_check check (repeat > 0)
);
