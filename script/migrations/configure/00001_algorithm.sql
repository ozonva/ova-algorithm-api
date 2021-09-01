-- +goose Up
create table algorithms (
  id serial primary key,
  subject text not null,
  description text not null
);

-- +goose Down
drop table if exists algorithms;
