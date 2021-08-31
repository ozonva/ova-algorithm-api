-- +goose Up
create table algorithms (
  id serial primary key,
  subject text,
  description text
);

-- +goose Down
drop table if exists algorithms;
