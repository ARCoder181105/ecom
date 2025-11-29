package orders

import (
	"database/sql"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := database.New(db)

	r.Group(func(protected chi.Router) {
		r.Use(utils.AuthMiddleware)

		r.Get("/orders", func(w http.ResponseWriter, r *http.Request) {
			handleUserOrders(w, r, q)
		})

	})

	return r

}
