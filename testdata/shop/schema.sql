-- A toy shop, the fixture for models' generators (roadmap M1.2). Its names
-- exercise the naming rules of ADR 0001 (D4, D7, D9). The app's table prefix
-- is shop_.

-- Plain values.
CREATE TYPE order_status AS ENUM ('pending', 'paid', 'shipped');

-- sqlc's character replacement: -, / and : become _, anything else that
-- isn't an ASCII letter or digit goes. 'pay-pal' and 'pay_pal' replace to
-- the same name, so sqlc numbers the second (value_<index>). '!!' replaces
-- to nothing and is numbered too; 'value_5' then collides with the number
-- '!!' took, and is numbered again. 'café' loses its é. The empty value
-- (PostgreSQL allows one) is numbered like '!!'.
CREATE TYPE payment_method AS ENUM (
  'credit-card', 'bank/transfer', 'cash:on:delivery',
  'pay-pal', 'pay_pal', '!!', 'value_5', 'café', ''
);

-- The initialism id, in the type's name and in values.
CREATE TYPE id_kind AS ENUM ('id', 'user_id', 'uuid');

-- No values at all (PostgreSQL allows it): no labels, and a Label() that
-- always returns nil.
CREATE TYPE no_values AS ENUM ();

-- A name that starts with a digit gets a leading underscore in Go; digits
-- inside names.
CREATE TYPE "2fa_method" AS ENUM ('totp', 'sms', 'u2f', 'level_3');

CREATE TABLE shop_customer (
  id bigint PRIMARY KEY,
  email text NOT NULL,
  id_kind id_kind NOT NULL,
  second_factor "2fa_method"
);
COMMENT ON TABLE shop_customer IS 'Customer';
COMMENT ON COLUMN shop_customer.email IS '  Where receipts go  ';

CREATE TABLE shop_order (
  id bigint PRIMARY KEY,
  customer_id bigint NOT NULL REFERENCES shop_customer (id),
  status order_status NOT NULL,
  payment payment_method,
  -- Starts with the prefix: stripped like a table name (zp does it too).
  shop_ref text NOT NULL,
  address_line2 text
);
-- White space only: the English text falls back to the readable name.
COMMENT ON COLUMN shop_order.shop_ref IS '   ';

CREATE TABLE shop_order_line (
  order_id bigint NOT NULL REFERENCES shop_order (id),
  line_no integer NOT NULL,
  title text NOT NULL,
  PRIMARY KEY (order_id, line_no)
);

-- No prefix: nothing to strip.
CREATE TABLE audit_log (
  id bigint PRIMARY KEY,
  happened_at timestamptz NOT NULL,
  what text NOT NULL
);

-- A view is listed with the tables.
CREATE VIEW shop_order_summary AS
  SELECT o.id, o.status, count(l.line_no) AS lines
  FROM shop_order o LEFT JOIN shop_order_line l ON l.order_id = o.id
  GROUP BY o.id, o.status;
COMMENT ON VIEW shop_order_summary IS 'Orders at a glance';

-- Listings (interface §4.3): every relation whose name ends in _view gets
-- a List<Type> method. A view with a nullable enum, a column named like a
-- keyword (order) and one named id, which sqlc spells ID.
CREATE VIEW shop_order_list_view AS
  SELECT o.id, o.status, o.payment, c.email, o.shop_ref AS "order"
  FROM shop_order o JOIN shop_customer c ON c.id = o.customer_id;

-- A table whose name ends in _view gets one too.
CREATE TABLE shop_page_view (
  id bigint PRIMARY KEY,
  path text NOT NULL,
  seen_at timestamptz NOT NULL
);

-- Column sets (interface §4.6, A4): a column of every kind, nullable and
-- not, and three of types no kind takes (an array, jsonb, interval), which
-- the set leaves out. gentest/shopptr reads the same columns through
-- sqlc's pointers and overrides (A5).
CREATE TABLE shop_kinds_view (
  id bigint PRIMARY KEY,
  name text NOT NULL,
  note varchar(40),
  ref uuid,
  small smallint NOT NULL,
  medium integer,
  ratio real,
  score double precision NOT NULL,
  price numeric(10, 2),
  active boolean,
  day date NOT NULL,
  at timestamp,
  seen timestamptz,
  span interval,
  state order_status,
  tags text[],
  meta jsonb
);
