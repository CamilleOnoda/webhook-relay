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

-- name: UpdateDelivery :exec
UPDATE deliveries
SET 
    status = $2,
    status_code = $3,
    response_body = $4,
    error_message = $5,
    delivery_duration_ms = $6,
    attempted_at = NOW(),
    attempt_count = attempt_count + 1,
    delivered_at = CASE
        WHEN $2 = 'success' THEN NOW()
        ELSE delivered_at
END
WHERE id = $1;

-- name: ListDeliveries :many
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
ORDER BY deliveries.created_at DESC;
