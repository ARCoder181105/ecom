package orders

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func handlePlaceOrder(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid user id"))
		return
	}

	var cartItems []mytypes.CreateOrderPayload
	if err := utils.ParseJson(r, &cartItems); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if len(cartItems) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("cart is empty"))
		return
	}

	// If any operation is undone, Everything is undone
	tx, err := db.Begin()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to start transaction"))
		return
	}

	defer tx.Rollback()

	qtx := database.New(db).WithTx(tx)

	var totalPrice = decimal.NewFromInt(0)

	productCache := make(map[uuid.UUID]database.Product)

	for _, item := range cartItems {
		prodID, err := uuid.Parse(item.ProductID)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid product id: %s", item.ProductID))
			return
		}

		product, err := qtx.GetProductByID(context.Background(), prodID)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("product not found: %s", item.ProductID))
			return
		}

		if int(product.StockQuantity) < item.Quantity {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("out of stock: %s", product.Name))
			return
		}

		itemTotal := product.Price.Mul(decimal.NewFromInt(int64(item.Quantity))) //price * quantity
		totalPrice = totalPrice.Add(itemTotal) // total+=price

		productCache[prodID] = product
	}

	order, err := qtx.CreateOrder(context.Background(), database.CreateOrderParams{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     "pending",
	})
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to create order"))
		return
	}

	for _, item := range cartItems {
		prodID, _ := uuid.Parse(item.ProductID) 
		product := productCache[prodID]

		_, err := qtx.CreateOrderItem(context.Background(), database.CreateOrderItemParams{
			OrderID:   order.ID,
			ProductID: prodID,
			Quantity:  int32(item.Quantity),
			Price:     product.Price, 
		})
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to create order item"))
			return
		}

		newStock := product.StockQuantity - int32(item.Quantity)

		_, err = qtx.UpdateProduct(context.Background(), database.UpdateProductParams{
			ID:            product.ID,
			Name:          product.Name,
			Description:   product.Description,
			Image:         product.Image,
			Price:         product.Price,
			StockQuantity: newStock,
			UserID:        product.UserID, 
		})

		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to update stock"))
			return
		}
	}

	if err := tx.Commit(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("failed to commit transaction"))
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Order placed successfully",
		"order_id": order.ID,
	})
}
