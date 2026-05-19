-- name: CreateEvent :one
INSERT INTO webhook_events(endpoint_id, event_type, payload, headers) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: ListEvents :many
SELECT * FROM webhook_events
ORDER BY received_at DESC;