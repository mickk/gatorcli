-- name: CreateFeed :one
INSERT INTO feeds (id, url, user_id, name, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;
-- name: GetFeeds :many
SELECT
  f.name,
  f.url,
  u.name as user_name
FROM feeds f
INNER JOIN users u
  ON f.user_id = u.id;
