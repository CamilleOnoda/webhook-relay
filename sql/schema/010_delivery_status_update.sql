-- +goose Up
UPDATE deliveries
SET status = 'retry_scheduled'
WHERE status = 'failed';

-- +goose Down
SELECT 1;
