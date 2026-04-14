create table if not exists tasks (
    id serial primary key,
    title varchar(255) not null
);