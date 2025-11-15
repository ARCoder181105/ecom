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
