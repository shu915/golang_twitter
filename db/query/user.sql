-- name: CreateUser :one
INSERT INTO users (email, password, token)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByToken :one
SELECT * FROM users WHERE token = $1;

-- name: UpdateUserIsActive :exec
UPDATE users SET is_active = $1,token = NULL WHERE token = $2;