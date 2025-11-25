-- +migrate Up

BEGIN;

CREATE TYPE IF NOT EXISTS company_type AS ENUM ('online_store', 'physical_store', 'multichannel_store');

CREATE TABLE IF NOT EXISTS companies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    company_type company_type NOT NULL,
    code CHAR(3) UNIQUE NOT NULL,
    contact_info JSONB DEFAULT '{}'::jsonb,
    settings JSONB DEFAULT '{}'::jsonb,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());

COMMIT;