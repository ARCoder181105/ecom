-- +goose Up
-- Create an ENUM type for safety, or just use VARCHAR with a check constraint
CREATE TYPE user_role AS ENUM ('customer', 'seller', 'admin');

ALTER TABLE users 
ADD COLUMN role user_role NOT NULL DEFAULT 'customer';

-- +goose Down
ALTER TABLE users DROP COLUMN role;
DROP TYPE user_role;