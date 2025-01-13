-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users (
                       uuid uuid default gen_random_uuid() primary key,
                       phone varchar not null,
                       password_hash varchar not null,
                       created_at timestamptz not null default clock_timestamp(),
                       updated_at timestamptz,
                       last_time_sms timestamptz not null default clock_timestamp(),
                       verification bool default False,
                       sms varchar default ''
);

create unique index user_phone_uniq_idx ON users(phone);



-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop index if exists users_email_uniq_idx;

drop table if exists users;