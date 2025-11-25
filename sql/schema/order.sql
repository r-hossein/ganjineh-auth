
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
    ));

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
    CONSTRAINT valid_quantities CHECK (returned_quantity <= quantity));

