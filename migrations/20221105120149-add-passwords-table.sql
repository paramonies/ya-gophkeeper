-- +migrate Up
create table if not exists passwords
(
    id              uuid default gen_random_uuid(),
    user_id         uuid not null,
    password        text not null,
    data_version    integer not null default 0,
    meta            text,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint passwords_pk primary key (id),
    constraint passwords_user_id_fk foreign key (user_id) references users (id)
);
-- +migrate Down
drop table passwords;