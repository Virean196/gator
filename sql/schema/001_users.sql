-- +goose Up
CREATE TABLE users(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP not NULL,
    updated_at TIMESTAMP not NULL,
    name TEXT UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;