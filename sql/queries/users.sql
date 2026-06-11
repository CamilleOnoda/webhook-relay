-- name: CreateUser :one
INSERT INTO users(name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE is_admin = false
AND id = $1
RETURNING id;

-- name: IsUserAdmin :one
SELECT is_admin FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;