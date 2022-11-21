-- +migrate Up
create table if not exists users
(
    id              uuid default gen_random_uuid(),
    login           text not null,
    password_hash   text not null,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint users_pk primary key (id),
    constraint login unique (login)
);
-- +migrate Down
drop table users;