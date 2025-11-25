package products

import (
	"database/sql"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := database.New(db)

	// public routes
	r.Get("/getAllProducts", func(w http.ResponseWriter, r *http.Request) {
		handleGetAllProducts(w, r, q)
	})

	r.Post("/upload", handleImageUpload)

	r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
		handleCreateProduct(w, r, q)
	})

	return r
}
