package orders

import (
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/google/uuid"
)

func handleUserOrders(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	userId := claims.UserID

	userIdUUID, err := uuid.Parse(userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	orders, err := q.ListOrdersByUser(r.Context(), userIdUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)

}
