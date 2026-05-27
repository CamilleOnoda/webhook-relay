-- name: CreateUser :one
INSERT INTO users(name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users
WHERE is_admin = false;

-- name: IsUserAdmin :one
SELECT is_admin FROM users WHERE id = $1;