-- name: CreateFollowedFeed :one
INSERT INTO followed_feeds (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFollowedFeedsByUserID :many
SELECT * FROM followed_feeds WHERE user_id = $1;

-- name: DeleteFollowedFeed :exec
DELETE FROM followed_feeds WHERE id = $1 AND user_id = $2;