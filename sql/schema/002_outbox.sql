-- +goose Up
CREATE TABLE outbox(
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL, -- kafka topic
    aggregate_type TEXT NOT NULL, -- type of entity involved (here just 'test')
    aggregate_id TEXT NOT NULL, -- kafka key -> test.id
    payload JSONB NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    published BOOLEAN NOT NULL DEFAULT FALSE
);

-- +goose Down
DROP TABLE test;