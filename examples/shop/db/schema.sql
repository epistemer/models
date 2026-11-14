-- The example shop's schema. Tables carry the prefix shop_, which the
-- labels leave out (genx.yaml sets table_prefix).
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped', 'cancelled');

CREATE TABLE shop_customer (
  id bigint PRIMARY KEY,
  email text NOT NULL,
  full_name text NOT NULL
);
COMMENT ON COLUMN shop_customer.email IS 'Email address';

CREATE TABLE shop_order (
  id bigint PRIMARY KEY,
  customer_id bigint NOT NULL REFERENCES shop_customer (id),
  status order_status NOT NULL DEFAULT 'pending',
  placed_at timestamptz NOT NULL
);
