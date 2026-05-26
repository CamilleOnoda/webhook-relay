-- name: CreateUser :one
INSERT INTO users(name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;