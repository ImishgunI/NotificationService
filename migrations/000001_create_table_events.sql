CREATE TYPE IF NOT EXISTS eventStatus AS ENUM ('NONE', 'ACCEPTED', 'PROCESSING', 'DONE', 'REJECTED', 'FAILED');

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    key varchar(30) NOT NULL,
    status eventStatus DEFAULT 'NONE',
    payload json NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_key ON events (key);