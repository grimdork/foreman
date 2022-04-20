package main

const initialTables = `-- Initial tables

-- Scouts table.
create table if not exists scouts (
	hostname text unique not null primary key,
	port smallint not null default 443,
	-- interval in seconds
	interval integer not null default 60,
	-- last check time
	last_check timestamp not null default current_timestamp,
	-- First failure time
	failed timestamp not null default current_timestamp - interval '1 year',
	-- last check status
	status smallint not null default 0,
	-- assignee who acknowledged the error
	assignee text not null default '',
	-- time of acknowledgement
	assigned timestamp not null default current_timestamp
);

-- Canaries table.
create table if not exists canaries (
	name text unique not null primary key,
	-- interval in seconds
	interval integer not null default 60,
	-- last check time
	last_check timestamp not null default current_timestamp,
	-- First failure time (or time when the system started noticing, anyway)
	failed timestamp not null default current_timestamp - interval '1 year',
	-- last check status
	status smallint not null default 0,
	-- assignee who acknowledged the error
	assignee text not null default '',
	-- time of acknowledgement
	assigned timestamp not null default current_timestamp
);

-- Authentication keys
create table if not exists keys (
	name text primary key,
	key text not null,
	admin boolean not null default false
);

--
-- Seed initial data.
--

-- Set default system key.
insert into keys(name,key,admin)
	values ('foreman','potrzebie',true)
	on conflict do nothing;
`
