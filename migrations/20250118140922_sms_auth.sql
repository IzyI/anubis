-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE sms_auth
(
    sms_uuid    uuid                 default gen_random_uuid() primary key,
    user_uuid   uuid        not null REFERENCES users (uuid),
    phone       bigint      not null,
    sms_code    varchar     not null,
    sms_service varchar     not null,
    id_send     varchar     not null,
    created_at  timestamptz not null default clock_timestamp()
);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table if exists sms_auth;
