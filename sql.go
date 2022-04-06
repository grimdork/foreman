package main

const initialTables = `-- Initial tables

-- Clients table.
create table if not exists clients (
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
	-- acknowledgement of error status
	acknowledgement boolean not null default false,
	-- assignee of the acknowledged error
	assignee text not null default ''
);

-- Authentication keys
create table if not exists keys (
	name text primary key,
	key text not null,
	expiry timestamp not null default now() + interval '1 month',
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
