// This file is user/route.go
package user

import (
	"database/sql"
	"net/http"

	db "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/go-chi/chi/v5"
)

// Routes sets up all user-related API endpoints.
func Routes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := db.New(database)

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, q)
	})
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, q)
	})

	return r
}

func handleLogin(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	// TODO: Implement login logic
	utils.RespondWithJSON(w, http.StatusNotImplemented, map[string]string{"message": "Login not yet implemented"})
}
