package orders

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func handleUserOrdersList(w http.ResponseWriter, r *http.Request, q *database.Queries) {
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
		totalPrice = totalPrice.Add(itemTotal)                                   // total+=price

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

func handleGetOrderById(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	orderIDStr := chi.URLParam(r, "orderID")
	if orderIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("order id is required"))
		return
	}

	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"))
		return
	}

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

	// check weather there is a order present
	order, err := q.GetOrderByID(r.Context(), database.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("order not found"))
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// after checking the order is presend we will fetch the order items
	items, err := q.GetOrderItems(r.Context(), orderID)
	if err != nil {
		if err == sql.ErrNoRows {
			items = []database.GetOrderItemsRow{}
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err)
			return
		}
	}

	response := map[string]interface{}{
		"order_id":    order.ID,
		"status":      order.Status,
		"total_price": order.TotalPrice,
		"created_at":  order.CreatedAt,
		"items":       items,
	}

	utils.RespondWithJSON(w, http.StatusOK, response)
}

func handleAdminUpdateStatus(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	if claims.Role != "admin" {
		utils.RespondWithError(w, http.StatusForbidden, fmt.Errorf("user is not admin"))
		return
	}

	var statusPayload mytypes.AdminUpdateStatusPayload
	if err := utils.ParseJson(r, &statusPayload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	if statusPayload.OrderID == "" || statusPayload.Status == "" {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("order_id and status are required"))
		return
	}

	orderIdUUID, err := uuid.Parse(statusPayload.OrderID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid order id"))
		return
	}

	// Attempt to update the order status and handle possible errors.
	if err := q.UpdateOrderStatus(r.Context(), database.UpdateOrderStatusParams{
		ID:     orderIdUUID,
		Status: statusPayload.Status,
	}); err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("order not found"))
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "order status updated successfully",
		"order_id": orderIdUUID,
		"status":   statusPayload.Status,
	})
}
