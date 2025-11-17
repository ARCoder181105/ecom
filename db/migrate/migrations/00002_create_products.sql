-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  image VARCHAR(255),
  price DECIMAL(10, 2) NOT NULL CHECK (price >= 0), -- Store price as a decimal
  stock_quantity INT NOT NULL DEFAULT 0 CHECK (stock_quantity >= 0),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products;
-- +goose StatementEnd