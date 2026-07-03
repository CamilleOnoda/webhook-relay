-- name: CreateUser :one
INSERT INTO users(name, email, hashed_password)
VALUES ($1, $2, $3)
RETURNING *;

-- name: DeleteUserByID :exec
DELETE FROM users
WHERE is_admin = false
AND id = $1
RETURNING id;

-- name: IsUserAdmin :one
SELECT is_admin FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserRecentActvity :many
SELECT
    e.id AS event_id,
    e.received_at,
    e.event_type,
    ep.name AS endpoint_name,
    COALESCE(d.status, 'pending') AS latest_delivery_status,
    d.status_code AS latest_delivery_status_code,
    COALESCE(d.attempt_count, 0) AS attempt_count,
    d.id AS delivery_id
FROM webhook_events e
JOIN webhook_endpoints ep
    ON ep.id = e.endpoint_id
LEFT JOIN LATERAL (
    SELECT
        status,
        status_code,
        created_at,
        attempt_count,
        id
    FROM deliveries
    WHERE deliveries.event_id = e.id
    ORDER BY created_at DESC
    LIMIT 1
) d ON true
WHERE ep.user_id = $1
ORDER BY e.received_at DESC
LIMIT 10;