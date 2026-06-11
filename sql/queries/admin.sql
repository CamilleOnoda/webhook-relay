-- name: GetAdminStats :one
SELECT
  (SELECT COUNT(*) FROM users) AS users,
  (SELECT COUNT(*) FROM webhook_events) AS events_received,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'success') AS successful_deliveries,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'failed') AS failed_deliveries;

  -- name: GetAllUsers :many
SELECT id, name, email, created_at, updated_at, is_admin FROM users;

-- name: GetAllEvents :many
SELECT * FROM webhook_events;

-- name: DeleteAllEndpoints :exec
DELETE FROM webhook_endpoints;