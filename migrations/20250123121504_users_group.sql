-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE users_group
(
    user_uuid  uuid    not null REFERENCES users (uuid),
    group_name varchar not null,
    service    varchar not null,
    CONSTRAINT users_group_pk PRIMARY KEY (user_uuid, group_name)
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table if exists users_group;
drop index if exists users_group_pk;