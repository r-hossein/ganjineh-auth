
  CREATE TYPE user_status AS ENUM ('active', 'inactive', 'suspended');
  CREATE TYPE session_status_type AS ENUM ('active','revoked','updated')
  CREATE TYPE order_status AS ENUM ('draft', 'pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded');
  CREATE TYPE order_from_type AS ENUM ('IN_PERSON', 'INSTAGRAM', 'WEBCITE', 'TELEGRAM');
  CREATE TYPE payment_status AS ENUM ('pending', 'paid', 'failed', 'refunded', 'partially_refunded');
  CREATE TYPE payment_method_type AS ENUM ('cash', 'ONLINE', 'cart');
  CREATE TYPE shipping_method_type AS ENUM ('IN_PERSON', 'POST', 'INTRA_CITY');
  CREATE TYPE company_type AS ENUM ('online_store', 'physical_store', 'multichannel_store');
  CREATE TYPE return_status_type AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');
  CREATE TYPE warehouse_type AS ENUM ('ONLINE', 'SOTORE', 'WAREHOUSE');
  CREATE TYPE inventory_status AS ENUM ('in_stock', 'out_of_stock', 'discontinued', 'pre_order');
  CREATE TYPE address_type AS ENUM ('billing', 'shipping', 'both');
  CREATE TYPE gender_type AS ENUM ('male', 'female', 'other', 'not_specified');
  CREATE TYPE product_status AS ENUM ('active', 'inactive', 'draft', 'archived');
  CREATE TYPE value_tier_type AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');
  CREATE TYPE behavior_segment_type AS ENUM ('requested', 'approved', 'rejected', 'completed', 'cancelled');
  CREATE TYPE inventory_change_type AS ENUM ('STOCK_IN', 'STOCK_OUT', 'ADJUSTMENT', 'CORRECTION');


CREATE OR REPLACE FUNCTION set_updated_at()
  RETURNS TRIGGER LANGUAGE plpgsql AS $$
  BEGIN
    IF NEW.last_login IS DISTINCT FROM OLD.last_login 
       AND (SELECT COUNT(*) FROM (
            SELECT key, value FROM jsonb_each(to_jsonb(NEW)) 
            WHERE key NOT IN ('last_login', 'updated_at')
            EXCEPT 
            SELECT key, value FROM jsonb_each(to_jsonb(OLD)) 
            WHERE key NOT IN ('last_login', 'updated_at')
        ) AS changes) = 0 THEN
        NEW.updated_at = OLD.updated_at;
    ELSE
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
  END;
$$;

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    permission_codes TEXT[] NOT NULL DEFAULT '{}',
    is_system_role BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW());

CREATE TABLE users (
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

CREATE TABLE user_sessions (
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

CREATE TABLE companies (
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

CREATE TABLE user_companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
    assigned_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    UNIQUE(user_id, company_id));


  CREATE TABLE brands (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    slug VARCHAR(300) NOT NULL,
    logo_url VARCHAR(500),
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
  );

  CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(300) NOT NULL,
    parent_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    image_url VARCHAR(500),
    brands_id TEXT[] DEFAULT '{}',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT no_self_reference CHECK (id IS NULL OR id <> parent_id)
  );

  CREATE TABLE attributes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,                      
    code VARCHAR(100) UNIQUE NOT NULL,               
    data_type VARCHAR(20) NOT NULL CHECK (data_type IN ('text', 'number', 'boolean', 'enum')),
    unit VARCHAR(50),
    code CHAR(3) UNIQUE,               
    options TEXT[],
    is_filterable BOOLEAN DEFAULT TRUE,
    is_variant BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
  );

  CREATE TABLE category_attributes (
    category_id INT NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    attribute_id INT NOT NULL REFERENCES attributes(id) ON DELETE RESTRICT,
    is_required BOOLEAN DEFAULT FALSE,
    sort_order INT DEFAULT 0,
    PRIMARY KEY (category_id, attribute_id)
  );

  CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    company_id INT NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    brand_id INT REFERENCES brands(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(300) NOT NULL,
    code CHAR(3),
    main_image_url VARCHAR(500),
    gallery_urls TEXT[] DEFAULT '{}',
    category_ids INTEGER[] DEFAULT '{}',
    base_price NUMERIC(14,2) CHECK (base_price >= 0),
    status product_status DEFAULT 'active',
    track_inventory BOOLEAN DEFAULT TRUE,
    min_stock_level INT DEFAULT 0,
    attributes JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (company_id, slug)
    UNIQUE (company_id, code)
  );

  CREATE TABLE product_variants (
    id SERIAL PRIMARY KEY,
    product_id INT NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    sku VARCHAR(100) UNIQUE NOT NULL,
    barcode VARCHAR(50),
    is_external_barcode BOOLEAN DEFAULT FALSE,
    price NUMERIC(14,2) CHECK (price >= 0),
    stock_quantity INT DEFAULT 0 CHECK (stock_quantity >= 0),
    is_active BOOLEAN DEFAULT TRUE,
    attributes JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
  );

  CREATE TABLE customers (
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
    UNIQUE(user_id, company_id)
  );

  CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE RESTRICT,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,
    agent_id INTEGER REFERENCES user_companies(id) ON DELETE SET NULL,
    order_number VARCHAR(9) NOT NULL,
    order_from order_from_type NOT NULL DEFAULT 'ONLINE',
    order_date TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    required_date TIMESTAMPTZ,
    shipped_date TIMESTAMPTZ,
    total_amount NUMERIC(14,2) NOT NULL CHECK (total_amount >= 0),
    discount_amount NUMERIC(14,2) DEFAULT 0 CHECK (discount_amount >= 0),
    shipping_amount NUMERIC(14,2) DEFAULT 0 CHECK (shipping_amount >= 0),
    final_amount NUMERIC(14,2) GENERATED ALWAYS AS (
      total_amount - discount_amount + shipping_amount
    ) STORED,
    status order_status DEFAULT 'pending' NOT NULL,
    payment_status payment_status DEFAULT 'pending' NOT NULL,
    payment_method payment_method_type DEFAULT 'cart' NOT NULL,
    payment_reference VARCHAR(50),
    shipping_method VARCHAR(100),
    tracking_number VARCHAR(200),
    shipping_address JSONB NOT NULL,
    customer_notes TEXT,
    internal_notes TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (company_id, order_number),
    CONSTRAINT valid_order_dates CHECK (
      (shipped_date IS NULL OR shipped_date >= order_date) AND
      (required_date IS NULL OR required_date >= order_date)
    )
  );

  CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE RESTRICT,
    product_id INTEGER REFERENCES products(id) ON DELETE SET NULL,
    product_name VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100) NOT NULL,
    product_data JSONB DEFAULT '{}'::jsonb,
    unit_price NUMERIC(12,2) NOT NULL CHECK (unit_price >= 0),
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    discount_percent NUMERIC(5,2) DEFAULT 0 CHECK (discount_percent BETWEEN 0 AND 100),
    discount_amount NUMERIC(12,2) DEFAULT 0 CHECK (discount_amount >= 0),
    line_total NUMERIC(14,2) GENERATED ALWAYS AS (
      (unit_price * quantity) - discount_amount
    ) STORED,
    warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE SET NULL,
    return_status return_status_type,
    returned_quantity INTEGER DEFAULT 0 CHECK (returned_quantity >= 0 AND returned_quantity <= quantity),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT valid_quantities CHECK (returned_quantity <= quantity)
  );

  CREATE TABLE warehouses (
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

  CREATE TABLE inventory (
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

  CREATE TABLE inventory_variant (
    variant_id INT NOT NULL REFERENCES product_variants(id) ON DELETE RESTRICT,
    warehouse_id INT NOT NULL REFERENCES warehouses(id) ON DELETE RESTRICT,
    current_stock INT DEFAULT 0 CHECK (current_stock >= 0),
    reserved_stock INT DEFAULT 0 CHECK (reserved_stock >= 0),
    available_stock INT GENERATED ALWAYS AS (current_stock - reserved_stock) STORED,
    last_updated TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (variant_id, warehouse_id)
  );


  CREATE TABLE inventory_adjustments (
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
    created_at TIMESTAMPTZ DEFAULT NOW()
  );

  CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_roles_updated_at BEFORE UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_brands_updated_at BEFORE UPDATE ON brands FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_categories_updated_at BEFORE UPDATE ON categories FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_warehouses_updated_at BEFORE UPDATE ON warehouses FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_products_updated_at BEFORE UPDATE ON products FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_customers_updated_at BEFORE UPDATE ON customers FOR EACH ROW EXECUTE FUNCTION set_updated_at();
  CREATE TRIGGER trg_orders_updated_at BEFORE UPDATE ON orders FOR EACH ROW EXECUTE FUNCTION set_updated_at();

  CREATE INDEX idx_users_phone ON users(phone_number);
  CREATE INDEX idx_companies_owner ON companies(owner_id);
  CREATE INDEX idx_products_company_slug ON products(company_id, slug);
  CREATE INDEX idx_products_company_sku ON products(company_id, sku);
  CREATE INDEX idx_inventory_product_warehouse ON inventory(product_id, warehouse_id);
  CREATE INDEX idx_inventory_available_partial ON inventory(product_id) WHERE (available_stock > 0);
  CREATE INDEX idx_orders_company_date ON orders(company_id, order_date DESC);
  CREATE INDEX idx_orders_number ON orders(order_number);
  CREATE INDEX idx_order_items_order ON order_items(order_id);
  CREATE INDEX idx_order_items_product ON order_items(product_id);

  CREATE OR REPLACE VIEW vw_product_total_stock AS
  SELECT p.company_id, p.id AS product_id, p.sku, p.name,
        SUM(i.current_stock) AS total_current_stock,
        SUM(i.reserved_stock) AS total_reserved_stock,
        SUM(i.current_stock) - SUM(i.reserved_stock) AS total_available_stock
  FROM products p
  LEFT JOIN inventory i ON i.product_id = p.id
  GROUP BY p.company_id, p.id, p.sku, p.name;
