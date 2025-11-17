package products

import (
	"context"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/ARCoder181105/ecom/utils"
)

func handleGetAllProducts(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	data, err := q.ListProducts(context.Background())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	

}
