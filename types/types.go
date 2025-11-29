package mytypes

import "time"

type UserResponse struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	PassWord  string `json:"password"`
}

type LoginUserPayload struct {
	Email    string `json:"email"`
	PassWord string `json:"password"`
}

type ProductResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Image         string    `json:"image"`
	Price         string    `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UserID        string    `json:"user_id"`
}

type CreateProductPayload struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	Price         string `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}

type CreateOrderPayload struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}
