-- name: CreatePost :exec
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
    $1,$2,$3,$4,$5,$6,$7,$8
);

-- name: GetPostsForUser :many
SELECT * FROM posts
INNER JOIN feeds ON feeds.id = posts.feed_id
INNER JOIN users ON feeds.user_id = users.id
WHERE users.id = $1;