-- name: CreatePost :one
INSERT INTO posts (user_id, content) VALUES ($1, $2)
RETURNING *;

-- name: GetPosts :many
SELECT * FROM posts;