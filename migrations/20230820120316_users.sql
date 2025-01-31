-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users
(
    uuid        uuid         default gen_random_uuid() primary key,
    created_at    timestamptz not null default clock_timestamp(),
    global_role int not null default 1
);



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop table if exists users;