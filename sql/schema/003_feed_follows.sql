-- +goose Up
CREATE TABLE feed_follows(
    id INTEGER PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP,
    user_id INTEGER NOT NULL,
    feed_id INTEGER NOT NULL,
    CONSTRAINT unique_feed_follow UNIQUE (user_id, feed_id),
    CONSTRAINT fk_users_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_feeds_id FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;