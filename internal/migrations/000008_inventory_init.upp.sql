-- +migrate Up

BEGIN;


CREATE TABLE IF NOT EXISTS warehouses (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    warehouse_type warehouse_type NOT NULL DEFAULT 'SOTRE',
    is_sale_online BOOLEAN DEFAULT TRUE,
    settings JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (company_id, code)
  );

CREATE TABLE IF NOT EXISTS inventory (
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    current_stock INTEGER DEFAULT 0 CHECK (current_stock >= 0),
    reserved_stock INTEGER DEFAULT 0 CHECK (reserved_stock >= 0),
    available_stock INTEGER GENERATED ALWAYS AS (current_stock - reserved_stock) STORED,
    last_stock_take TIMESTAMPTZ,
    last_stock_take_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    low_stock_threshold INTEGER DEFAULT 10,
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    version INTEGER DEFAULT 1,
    PRIMARY KEY (product_id, warehouse_id)
  );

CREATE TABLE IF NOT EXISTS inventory_variant (
    variant_id INT NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT,
    warehouse_id INT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    current_stock INT DEFAULT 0 CHECK (current_stock >= 0),
    reserved_stock INT DEFAULT 0 CHECK (reserved_stock >= 0),
    available_stock INT GENERATED ALWAYS AS (current_stock - reserved_stock) STORED,
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (variant_id, warehouse_id)
  );


CREATE TABLE IF NOT EXISTS inventory_adjustments (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    warehouse_id INTEGER NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    change_by INTEGER REFERENCES users(id) ON DELETE SET NULL,
    change_type inventory_change_type NOT NULL,
    quantity_change INTEGER NOT NULL,
    previous_stock INTEGER,
    new_stock INTEGER,
    reason TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW());

COMMIT;