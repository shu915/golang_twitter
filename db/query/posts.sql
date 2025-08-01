-- name: CreatePost :one
INSERT INTO posts (user_id, content) VALUES ($1, $2)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: CountPosts :one
SELECT COUNT(*) FROM posts;