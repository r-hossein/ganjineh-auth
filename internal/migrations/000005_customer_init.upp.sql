-- +migrate Up

BEGIN;

CREATE TYPE value_tier_type AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');

CREATE TYPE behavior_segment_type AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    agent_id INTEGER REFERENCES user_companies(id),
    gender gender_type DEFAULT 'not_specified',
    date_of_birth DATE,
    is_spam BOOLEAN DEFAULT FALSE,
    value_tier value_tier_type DEFAULT 'standard',
    behavior_segment behavior_segment_type DEFAULT 'new',
    total_orders INTEGER DEFAULT 0,
    total_spent NUMERIC(14,2) DEFAULT 0,
    average_order_value NUMERIC(12,2) DEFAULT 0,
    registration_date TIMESTAMPTZ DEFAULT NOW(),
    first_purchase_date TIMESTAMPTZ,
    last_purchase_date TIMESTAMPTZ,
    addresses JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, company_id));

COMMIT;