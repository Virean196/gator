-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS(
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES(
        $1,$2,$3,$4,$5
    )
    RETURNING *
) SELECT inserted_feed_follow.*,
    feeds.name AS feed_name,
    users.name AS user_name
    FROM inserted_feed_follow
    INNER JOIN users on inserted_feed_follow.user_id = users.id
    INNER JOIN feeds on inserted_feed_follow.feed_id = feeds.id;

-- name: GetFollowing :many
SELECT feeds.name FROM feeds
INNER JOIN feed_follows ON feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
USING users, feeds
WHERE feed_follows.user_id = users.id
AND feed_follows.feed_id = feeds.id
AND users.name = $1
AND feeds.url = $2;