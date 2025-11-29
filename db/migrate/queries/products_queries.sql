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
WHERE 
    (name ILIKE '%' || $3 || '%' OR description ILIKE '%' || $3 || '%') -- Search logic
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

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

-- name: DeleteProduct :one
-- Used by Sellers/Customers: Deletes only if they own it
DELETE FROM products
WHERE id = $1 AND user_id = $2
RETURNING id;

-- name: DeleteProductByAdmin :one
-- Used by Admins: Deletes by ID (ignores ownership)
DELETE FROM products
WHERE id = $1
RETURNING id;

-- name: CountProducts :one
SELECT COUNT(*) FROM products
WHERE 
    (name ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%');