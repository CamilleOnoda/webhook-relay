-- name: CreateEndpoint :one
INSERT INTO webhook_endpoints (name, target_url, description, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetEndpointByID :one
SELECT * FROM webhook_endpoints
WHERE id = $1;

-- name: GetEndpointByIDAndUserID :one
SELECT * FROM webhook_endpoints
WHERE id = $1
AND user_id = $2;

-- name: DeleteAllEndpoints :exec
DELETE FROM webhook_endpoints;

-- name: DeleteEndpointByIDAndUserID :exec
DELETE FROM webhook_endpoints
WHERE id = $1
AND user_id = $2;

-- name: GetEndpointsByUserID :many
SELECT * FROM webhook_endpoints
WHERE user_id = $1
ORDER BY created_at DESC;