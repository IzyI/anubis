-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table phone_auth
(
    phone         bigint      not null,
    country_code  int         not null,
    password_hash varchar     not null,
    created_at    timestamptz not null default clock_timestamp(),
    verification  bool        not null default False,
    user_uuid     uuid        not null REFERENCES users (uuid)
);

CREATE UNIQUE INDEX user_phone_uniq_idx ON phone_auth (phone);



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop table if exists phone_auth;
drop index if exists user_phone_uniq_idx;