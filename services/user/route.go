// This file is user/route.go
package user

import (
	"database/sql"
	"net/http"

	db "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	// "github.com/ARCoder181105/ecom/utils"
	"github.com/go-chi/chi/v5"
)

// Routes sets up all user-related API endpoints.
func Routes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := db.New(database)

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
