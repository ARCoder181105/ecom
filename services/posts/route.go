package posts

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
)

func Routes(db *sql.DB) chi.Router {
	r := chi.NewRouter()

	return r
}
