-- name: CreateOrder :one
INSERT INTO orders (user_id, total_price, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, product_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders 
WHERE id = $1 AND user_id = $2 
LIMIT 1;

-- name: ListOrdersByUser :many
SELECT * FROM orders 
WHERE user_id = $1 
ORDER BY created_at DESC;