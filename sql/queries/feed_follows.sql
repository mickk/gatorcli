-- name: CreateFeedFollow :one
WITH inserted as (
INSERT INTO feed_follows (id, user_id, feed_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
) 
SELECT
  i.*,
  f.name AS feed_name,
  u.name AS user_name
FROM inserted i
INNER JOIN feeds f
  ON i.feed_id = f.id
INNER JOIN users u
  ON i.user_id = u.id;
-- name: GetFeedFollowsForUser :many
SELECT
  ff.*,
  f.name as feed_name,
  u.name as user_name
FROM feed_follows ff
INNER JOIN users u
  ON ff.user_id = u.id
INNER JOIN feeds f
  ON ff.feed_id = f.id
WHERE u.name = $1;
