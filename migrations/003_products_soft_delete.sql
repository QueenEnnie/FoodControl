ALTER TABLE products
    ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_products_deleted_at
    ON products (deleted_at);

CREATE INDEX IF NOT EXISTS idx_products_active_lookup
    ON products (lower(name), lower(unit))
    WHERE deleted_at IS NULL;
