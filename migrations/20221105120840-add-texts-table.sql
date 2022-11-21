-- +migrate Up
create table if not exists texts
(
    id              uuid default gen_random_uuid(),
    user_id         uuid not null,
    text_data       text,
    data_version    integer not null default 0,
    meta            text,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint texts_pk primary key (id),
    constraint texts_user_id_fk foreign key (user_id) references users (id)
);
-- +migrate Down
drop table texts;
