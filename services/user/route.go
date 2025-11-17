// This file is user/route.go
package user

import (
	"database/sql"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/go-chi/chi/v5"
)

// Routes sets up all user-related API endpoints.
func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := database.New(db)

	// Public routes
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, q)
	})

	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, q)
	})

	// Protected routes
	// r.Group(func(pr chi.Router) {
	// 	pr.Use(utils.AuthMiddleware)

	// 	pr.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
	// 		handleProfile(w, r, q)
	// 	})

	// 	pr.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
	// 		handleOrders(w, r, q)
	// 	})
	// })

	return r
}
