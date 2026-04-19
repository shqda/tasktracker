-- +goose Up
create table tasks (
     id serial primary key,
     title varchar(255) not null
);

-- +goose Down
drop table if exists tasks;