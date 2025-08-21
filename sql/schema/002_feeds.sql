-- +goose Up
CREATE TABLE feeds(
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    user_id INTEGER NOT NULL,
    CONSTRAINT fk_users_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;