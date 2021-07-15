create table if not exists clients (
	id serial PRIMARY KEY,
	"host" text not null,
    "name" text not null,
	"key" text not null
);

create table if not exists users (
	id serial primary key,
	client_id int not null references clients (id),
	external_id text not null,
	"name" text
);

create table if not exists rooms (
	id serial primary key,
	"name" text,
	"key" text not null
);

create table if not exists roles (
    id serial primary key,
    "name" text not null
);

create table if not exists subs (
    user_id bigint not null references users (id),
    room_id bigint not null references rooms (id),
    role_id int not null references roles (id)
);