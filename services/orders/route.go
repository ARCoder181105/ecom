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
			handleUserOrdersList(w, r, q)
		})

		r.Get("/orders/{orderID}", func(w http.ResponseWriter, r *http.Request) {
			handleGetOrderById(w, r, q)
		})

		// Need Database transaction
		r.Post("/placeOrder", func(w http.ResponseWriter, r *http.Request) {
			handlePlaceOrder(w, r, db)
		})

		r.Post("/updateOrderStatus", func(w http.ResponseWriter, r *http.Request) {
			handleAdminUpdateStatus(w, r, q)
		})

	})

	return r

}
