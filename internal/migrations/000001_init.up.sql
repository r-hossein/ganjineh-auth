-- +migrate Up
-- Initial schema migration created from `sql/schema.sql`

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TYPE IF NOT EXISTS user_status AS ENUM ('active', 'inactive', 'suspended');
CREATE TYPE IF NOT EXISTS session_status_type AS ENUM ('active','revoked','updated');
CREATE TYPE IF NOT EXISTS company_type AS ENUM ('online_store', 'physical_store', 'multichannel_store');

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    permission_codes TEXT[] NOT NULL DEFAULT '{}',
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW());

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    no_id CHAR(10),
    phone_number CHAR(11) NOT NULL UNIQUE,
    email VARCHAR(100),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    profile_data JSONB DEFAULT '{}'::jsonb,
    password_hash  VARCHAR(256),
    status user_status NOT NULL DEFAULT 'active',
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    is_phone_verified BOOLEAN NOT NULL DEFAULT FALSE,
    is_password_set BOOLEAN NOT NULL DEFAULT FALSE,
    is_profile_complete BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW());

CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    device_info JSONB NOT NULL DEFAULT '{}'::jsonb,
    refresh_token_hash VARCHAR(128) NOT NULL,
    refresh_token_created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    refresh_token_expires_at TIMESTAMPTZ NOT NULL,
    last_active TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status session_status_type NOT NULL DEFAULT 'active',
    revoked_at TIMESTAMPTZ,
    revoked_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT valid_token_dates CHECK (refresh_token_expires_at > refresh_token_created_at));

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

CREATE TABLE IF NOT EXISTS user_companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    assigned_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, company_id));
