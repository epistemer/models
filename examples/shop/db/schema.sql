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

-- Its name ends in _view, so models' views generator gives it a listing
-- method, ListShopOrderListView.
CREATE VIEW shop_order_list_view AS
  SELECT o.id, o.status, o.placed_at, c.email
  FROM shop_order o JOIN shop_customer c ON c.id = o.customer_id;
COMMENT ON VIEW shop_order_list_view IS 'Orders';
