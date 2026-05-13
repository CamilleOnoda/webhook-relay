-- name: CreateEndpoint :one
INSERT INTO webhook_endpoints (name, target_url)
VALUES ($1, $2)
RETURNING *;