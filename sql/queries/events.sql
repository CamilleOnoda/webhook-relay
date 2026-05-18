-- name: CreateEvent :one
INSERT INTO webhook_events(endpoint_id, event_type, payload, headers) 
VALUES ($1, $2, $3, $4) 
RETURNING *;