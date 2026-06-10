-- name: GetAdminStats :one
SELECT
  (SELECT COUNT(*) FROM users) AS users,
  (SELECT COUNT(*) FROM webhook_events) AS events_received,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'success') AS successful_deliveries,
  (SELECT COUNT(*) FROM deliveries WHERE status = 'failed') AS failed_deliveries;