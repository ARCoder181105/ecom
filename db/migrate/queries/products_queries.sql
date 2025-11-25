-- name: CreateProduct :one
INSERT INTO products (
    name, description, image, price, stock_quantity, user_id
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products
WHERE id = $1
LIMIT 1;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY created_at DESC;

-- name: UpdateProduct :one
UPDATE products
SET 
    name = $2,
    description = $3,
    image = $4,
    price = $5,
    stock_quantity = $6
WHERE id = $1 AND user_id = $7
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND user_id = $2;
