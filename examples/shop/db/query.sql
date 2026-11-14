-- name: ListOrders :many
SELECT * FROM shop_order ORDER BY placed_at DESC;
