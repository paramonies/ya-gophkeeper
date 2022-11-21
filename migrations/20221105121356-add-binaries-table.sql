-- +migrate Up
create table if not exists binaries
(
    id              uuid default gen_random_uuid(),
    user_id         uuid not null,
    binarу_data     text,
    data_version    integer not null default 0,
    meta            text,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint binaries_pk primary key (id),
    constraint binaries_user_id_fk foreign key (user_id) references users (id)
);
-- +migrate Down
drop table binaries;
