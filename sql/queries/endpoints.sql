-- name: CreateEndpoint :one
INSERT INTO webhook_endpoints (name, target_url)
VALUES ($1, $2)
RETURNING *;

-- name: ListEndpoints :many
SELECT * FROM webhook_endpoints
ORDER BY created_at DESC;

-- name: GetEndpointByID :one
SELECT * FROM webhook_endpoints
WHERE id = $1;

-- name: DeleteAllEndpoints :exec
DELETE FROM webhook_endpoints;

-- name: DeleteEndpointByID :exec
DELETE FROM webhook_endpoints
WHERE id = $1;