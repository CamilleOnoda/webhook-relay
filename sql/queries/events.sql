-- name: CreateEvent :one
INSERT INTO webhook_events(endpoint_id, event_type, payload, headers) 
VALUES ($1, $2, $3, $4) 
RETURNING *;

-- name: ListEventsByUser :many
SELECT
    webhook_events.id,
    webhook_endpoints.name AS endpoint_name,
    webhook_events.event_type,
    webhook_events.received_at
FROM webhook_events
JOIN webhook_endpoints
    ON webhook_events.endpoint_id = webhook_endpoints.id
WHERE webhook_endpoints.user_id = $1
ORDER BY webhook_events.received_at DESC;

-- name: GetEventByID :one
SELECT * FROM webhook_events
WHERE id = $1;