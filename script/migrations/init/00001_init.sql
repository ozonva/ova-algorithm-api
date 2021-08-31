-- +goose NO TRANSACTION
-- +goose Up

create database ova;
create user melkozer with password 'melkozer';
grant all privileges on database ova to melkozer;

-- +goose Down

drop database if exists ova;
drop user if exists melkozer;
