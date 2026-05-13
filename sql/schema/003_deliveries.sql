-- +goose Up
CREATE TABLE deliveries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID NOT NULL REFERENCES webhook_events(id) ON DELETE RESTRICT,
    target_url TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    status_code INT,
    response_body TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    attempted_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    attempt_count INT NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ,
    delivery_duration_ms INT
);

-- +goose Down
DROP TABLE deliveries;