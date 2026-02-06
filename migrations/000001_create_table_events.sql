CREATE TYPE IF NOT EXISTS eventStatus AS ENUM ('ACCEPTED', 'PROCESSING', 'DONE', 'REJECTED', 'FAILED');

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_key UUID NOT NULL UNIQUE,
    status eventStatus,
    payload jsonb NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_key ON events (event_key);