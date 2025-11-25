-- +goose Up
ALTER TABLE products ADD COLUMN user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE products DROP COLUMN user_id;