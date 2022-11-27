-- +migrate Up
create table if not exists cards
(
    id              uuid default gen_random_uuid(),
    user_id         uuid not null,
    number          text not null,
    owner           text not null,
    expiration_date varchar(12) not null,
    cvv             integer not null,
    version         integer not null default 0,
    meta            text,
    created_at      timestamp default now(),
    deleted_at      timestamp,

    constraint cards_pk primary key (id),
    constraint cards_user_id_fk foreign key (user_id) references users (id)
);
-- +migrate Down
drop table cards;
