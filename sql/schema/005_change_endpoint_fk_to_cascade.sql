-- +goose Up
ALTER TABLE webhook_events
DROP CONSTRAINT webhook_events_endpoint_id_fkey;

ALTER TABLE webhook_events
ADD CONSTRAINT webhook_events_endpoint_id_fkey
FOREIGN KEY (endpoint_id)
REFERENCES webhook_endpoints(id)
ON DELETE CASCADE;

-- +goose Down
ALTER TABLE webhook_events
DROP CONSTRAINT webhook_events_endpoint_id_fkey;

ALTER TABLE webhook_events
ADD CONSTRAINT webhook_events_endpoint_id_fkey
FOREIGN KEY (endpoint_id)
REFERENCES webhook_endpoints(id)
ON DELETE RESTRICT;