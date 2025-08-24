-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds
where url = $1;

-- name: GetFeedUser :one
SELECT * FROM users
INNER JOIN feeds ON users.id = feeds.user_id
WHERE user_id = $1;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1
WHERE id = $2;

-- name: GetNextFeedTofetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;