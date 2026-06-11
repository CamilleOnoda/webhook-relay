-- name: GetAdminStats :one
SELECT
  (SELECT COUNT(*) FROM users) AS users,
  (SELECT COUNT(*) FROM webhook_events) AS events_received,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'success') AS successful_deliveries,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'failed') AS failed_deliveries;

-- name: GetAllUsers :many
SELECT 
  id, 
  name, 
  email, 
  created_at, 
  updated_at, 
  is_admin 
FROM users;

-- name: GetAllEvents :many
SELECT
  events.id,
  events.endpoint_id,
  events.event_type,
  events.received_at,
  ep.name AS endpoint_name
FROM webhook_events events
JOIN webhook_endpoints ep
  ON ep.id = events.endpoint_id
ORDER BY events.received_at DESC
LIMIT 10;

-- name: GetAllEndpoints :many
SELECT 
  name, 
  is_active, 
  created_at, 
  user_id 
FROM webhook_endpoints
ORDER BY created_at DESC
LIMIT 10;

-- name: GetAlldeliveries :many
SELECT 
  d.event_id, 
  d.status, 
  d.status_code, 
  d.target_url, 
  d.created_at,
  ep.name AS endpoint_name
FROM deliveries d
JOIN webhook_events e
  ON e.id = d.event_id
JOIN webhook_endpoints ep
  ON ep.id = e.endpoint_id
ORDER BY d.created_at DESC
LIMIT 10;

-- name: DeleteAllEndpoints :exec
DELETE FROM webhook_endpoints;