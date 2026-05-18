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
