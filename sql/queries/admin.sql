-- name: GetAdminStats :one
SELECT
  (SELECT COUNT(*) FROM users) AS users,
  (SELECT COUNT(*) FROM webhook_events) AS events_received,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'success') AS successful_deliveries,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'dead_letter') AS dead_letter,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'retry_scheduled') AS retry_scheduled_deliveries;

-- name: GetAllUsers :many
SELECT 
  id, 
  name, 
  email, 
  created_at, 
  updated_at, 
  is_admin 
FROM users;

-- name: GetRecentActivity :many
SELECT
  e.id AS event_id,
  e.received_at,
  ep.name AS endpoint_name,
  u.name AS user_name,
  e.event_type,
  COALESCE(d.status, 'pending') AS latest_delivery_status,
  d.status_code AS latest_delivery_status_code,
  COALESCE(d.attempt_count, 0) AS attempt_count,
  d.id AS delivery_id
FROM webhook_events e
JOIN webhook_endpoints ep
  ON ep.id = e.endpoint_id
JOIN users u
  ON u.id = ep.user_id
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
ORDER BY e.received_at DESC
LIMIT 10;

-- name: GetAllEvents :many
SELECT
  e.received_at,
  e.event_type,
  ep.name AS endpoint_name,
  u.name AS user_name
FROM webhook_events e
JOIN webhook_endpoints ep
  ON ep.id = e.endpoint_id
JOIN users u
  ON u.id = ep.user_id
ORDER BY e.received_at DESC
LIMIT 10;

-- name: GetAllEndpoints :many
SELECT 
  ep.name, 
  ep.is_active, 
  ep.created_at, 
  ep.user_id,
  u.name as user_name
FROM webhook_endpoints ep
JOIN users u
  ON u.id = ep.user_id
ORDER BY ep.created_at DESC
LIMIT 10;

-- name: GetAllDeliveries :many
SELECT
  d.created_at,
  d.status,
  d.status_code,
  ep.name AS endpoint_name,
  u.name AS user_name
FROM deliveries d
JOIN webhook_events e
  ON e.id = d.event_id
JOIN webhook_endpoints ep
  ON ep.id = e.endpoint_id
JOIN users u
  ON u.id = ep.user_id
ORDER BY d.created_at DESC
LIMIT 10;

-- name: DeleteAllEndpoints :exec
DELETE FROM webhook_endpoints;

-- name: GetAdminDeliveryDetails :one
SELECT
  d.id,
  ep.name AS endpoint_name,
  d.target_url,
  d.status,
  d.attempt_count,
  d.status_code,
  d.error_message,
  d.next_retry_at,
  d.created_at,
  d.attempted_at,
  d.delivered_at
FROM deliveries d
JOIN webhook_events e
  ON e.id = d.event_id
JOIN webhook_endpoints ep
  ON ep.id = e.endpoint_id
JOIN users u
  ON u.id = ep.user_id
WHERE d.id = $1;