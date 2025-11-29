-- +migrate Up
-- Initial schema migration created from `sql/schema.sql`

BEGIN;

CREATE TABLE IF NOT EXISTS errors (
    id BIGSERIAL PRIMARY KEY,
    http_code INT NOT NULL,
    status_code INT NOT NULL,
    message TEXT NOT NULL,
    stack_trace TEXT,
    endpoint VARCHAR(255),
    method VARCHAR(10),
    query_params JSONB NOT NULL DEFAULT '{}'::jsonb,
    request_body JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT NOW()
);

COMMIT;