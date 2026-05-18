-- +goose Up
ALTER TABLE webhook_endpoints
ADD COLUMN description TEXT;

-- +goose Down
ALTER TABLE webhook_endpoints
DROP COLUMN description;
