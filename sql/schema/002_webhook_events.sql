-- +goose Up
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    endpoint_id UUID NOT NULL REFERENCES webhook_endpoints(id) ON DELETE RESTRICT,
    event_type TEXT,
    payload JSONB NOT NULL,
    headers JSONB NOT NULL,
    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

--+goose Down
DROP TABLE webhook_events;