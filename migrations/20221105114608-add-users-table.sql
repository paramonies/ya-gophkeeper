-- +migrate Up
create table if not exists users
(
    id              uuid default gen_random_uuid(),
    email           text not null,
    password_hash   text not null,
    data_version    integer not null default 0,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint users_pk primary key (id),
    constraint email unique (email)
);
-- +migrate Down
drop table users;