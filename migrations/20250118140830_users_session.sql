-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE users_session
(
    ip            inet    not null,
    uuid          uuid             default gen_random_uuid() primary key,
    id_device     varchar not null,
    type          varchar not null,
    revoked       bool    not null default false,
    refresh_token varchar not null,
    user_uuid     uuid    not null REFERENCES users (uuid),
    id_service    int     not null
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop table if exists users_session;