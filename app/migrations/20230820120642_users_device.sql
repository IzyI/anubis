-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table users_devices
(
    id    serial primary key,
    ip inet  not null,
    id_device varchar not null,
    type varchar not null,
    user_uuid    uuid REFERENCES users (uuid)
);

create index id_device_idx ON users_devices (id_device);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd


drop index if exists id_device_idx;

drop table if exists users_device;