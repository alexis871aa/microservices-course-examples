-- +goose Up
-- создаем таблицу пользователей
CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   text,
    email      text,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- удаляем таблицу пользователей
drop table if exists users;
