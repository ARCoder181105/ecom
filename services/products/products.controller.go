package products

import (
	"context"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
)

func handleGetAllProducts(w http.ResponseWriter, _ *http.Request, q *database.Queries) {
	data, err := q.ListProducts(context.Background())
	if err != nil {
		http.Error(w, "Unable to list products", http.StatusInternalServerError)
		return
	}

	var responseProducts []mytypes.ProductResponse

	for _, row := range data {
		responseProducts = append(responseProducts, mytypes.ProductResponse{
			ID:            row.ID.String(),
			Name:          row.Name,
			Description:   row.Description,
			Price:         row.Price.String(), 
			StockQuantity: int(row.StockQuantity),
			CreatedAt:     row.CreatedAt,
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, responseProducts)
}
