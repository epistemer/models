-- name: GetOrder :one
SELECT * FROM shop_order WHERE id = $1;

-- name: ListOrderLines :many
SELECT * FROM shop_order_line WHERE order_id = $1 ORDER BY line_no;

-- name: OrderSummaries :many
SELECT * FROM shop_order_summary;
