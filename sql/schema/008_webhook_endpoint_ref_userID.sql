-- +goose Up
ALTER TABLE webhook_endpoints
ADD COLUMN user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE webhook_endpoints
DROP COLUMN user_id;