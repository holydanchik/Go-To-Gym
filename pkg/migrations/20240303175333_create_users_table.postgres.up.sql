create table users (
    id serial primary key,
    username varchar(255) not null,
    email varchar(255) not null
);