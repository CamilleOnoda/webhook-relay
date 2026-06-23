-- name: CreateDelivery :one
INSERT INTO deliveries (
    event_id,
    target_url,
    status,
    status_code,
    response_body,
    error_message
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListDeliveriesByUser :many
SELECT 
    deliveries.id,
    webhook_endpoints.name AS endpoint_name,
    deliveries.target_url,
    deliveries.status,
    deliveries.status_code,
    deliveries.created_at,
    deliveries.delivery_duration_ms
FROM deliveries
JOIN webhook_events ON deliveries.event_id = webhook_events.id
JOIN webhook_endpoints ON webhook_events.endpoint_id = webhook_endpoints.id
WHERE webhook_endpoints.user_id = $1
ORDER BY deliveries.created_at DESC;

-- name: GetReadyDeliveries :many
SELECT * FROM deliveries
WHERE status = 'pending' OR (status = 'retry_scheduled' AND next_retry_at <= NOW())
ORDER BY COALESCE(next_retry_at, created_at) ASC
LIMIT $1;

-- name: UpdateDeliveryState :exec
UPDATE deliveries
SET 
    status = $2,
    status_code = $3,
    response_body = $4,
    error_message = $5,
    attempt_count = attempt_count + 1,
    attempted_at = NOW(),
    next_retry_at = $6,
    delivered_at = $7,
    delivery_duration_ms = $8
WHERE id = $1;

-- name: ListDeadLetterDeliveries :many
SELECT 
    d.id,
    d.target_url,
    d.status,
    d.status_code,
    d.created_at,
    d.delivery_duration_ms,
    d.response_body,
    d.error_message,
    e.name AS endpoint_name,
    u.name AS user_name
FROM deliveries d
JOIN webhook_events ev ON d.event_id = ev.id
JOIN webhook_endpoints e ON ev.endpoint_id = e.id
JOIN users u ON e.user_id = u.id
WHERE d.status = 'dead_letter'
ORDER BY d.created_at DESC;


-- name: ReplayDeadLetterDelivery :exec
UPDATE deliveries
SET 
    status = 'pending',
    attempt_count = 0,
    next_retry_at = NULL,
    error_message = NULL,
    status_code = NULL,
    response_body = NULL
WHERE id = $1 AND status = 'dead_letter';
