-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd


CREATE TABLE users_info
(
    user_uuid uuid not null REFERENCES users (uuid),
    email     varchar,
    nickname  varchar,
    CONSTRAINT users_info_pk PRIMARY KEY (user_uuid, email)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop table if exists users_info;
drop index if exists users_info_pk;