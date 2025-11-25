package products

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

func handleGetAllProducts(w http.ResponseWriter, _ *http.Request, q *database.Queries) {
	data, err := q.ListProducts(context.Background())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, fmt.Errorf("unable to list products"))
		return
	}

	var responseProducts []mytypes.ProductResponse

	for _, row := range data {
		responseProducts = append(responseProducts, mytypes.ProductResponse{
			ID:            row.ID.String(),
			Name:          row.Name,
			Description:   row.Description,
			Image:         row.Image.String,
			Price:         row.Price.String(),
			StockQuantity: int(row.StockQuantity),
			CreatedAt:     row.CreatedAt,
			UserID:        row.UserID.String(),
		})
	}

	utils.RespondWithJSON(w, http.StatusOK, responseProducts)
}

func handleGetProductByID(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	productIDStr := chi.URLParam(r, "productID")

	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid product id format"))
		return
	}

	product, err := q.GetProductByID(context.Background(), productID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("product not found"))
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	resp := mytypes.ProductResponse{
		ID:            product.ID.String(),
		Name:          product.Name,
		Description:   product.Description,
		Image:         product.Image.String,
		Price:         product.Price.String(),
		StockQuantity: int(product.StockQuantity),
		CreatedAt:     product.CreatedAt,
		UserID:        product.UserID.String(),
	}

	utils.RespondWithJSON(w, http.StatusOK, resp)
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	var payload mytypes.CreateProductPayload

	// Parse JSON
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	// Extract user from JWT
	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	if claims.Role != "seller" && claims.Role != "admin" {
		utils.RespondWithError(w, http.StatusForbidden, fmt.Errorf("only sellers can create products"))
		return
	}

	userID, _ := uuid.Parse(claims.UserID)

	// Validate price
	price, err := decimal.NewFromString(payload.Price)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid price format"))
		return
	}

	// Insert product
	product, err := q.CreateProduct(context.Background(), database.CreateProductParams{
		Name:          payload.Name,
		Description:   payload.Description,
		Image:         sql.NullString{String: payload.Image, Valid: payload.Image != ""},
		Price:         price,
		StockQuantity: int32(payload.StockQuantity),
		UserID:        userID,
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, product)
}

func handleUpdateProduct(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid product id"))
		return
	}

	var payload mytypes.CreateProductPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid request body"))
		return
	}

	price, err := decimal.NewFromString(payload.Price)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid price"))
		return
	}

	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}
	userID, _ := uuid.Parse(claims.UserID)

	updatedProduct, err := q.UpdateProduct(context.Background(), database.UpdateProductParams{
		ID:            productID,
		Name:          payload.Name,
		Description:   payload.Description,
		Image:         sql.NullString{String: payload.Image, Valid: payload.Image != ""},
		Price:         price,
		StockQuantity: int32(payload.StockQuantity),
		UserID:        userID,
	})

	if err == sql.ErrNoRows {
		utils.RespondWithError(w, http.StatusForbidden, fmt.Errorf("you do not own this product"))
		return
	}

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedProduct)
}

func handleDeleteProduct(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)

	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid product ID"))
		return
	}

	// 1. Get User Claims (ID and Role)
	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
		return
	}

	// 2. Determine which Query to use based on Role
	if claims.Role == "admin" {
		// --- ADMIN PATH ---
		_, err = q.DeleteProductByAdmin(context.Background(), productID)

		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Errorf("product not found"))
			return
		}
	} else {
		// --- SELLER/USER PATH ---
		userID, _ := uuid.Parse(claims.UserID)
		_, err = q.DeleteProduct(context.Background(), database.DeleteProductParams{
			ID:     productID,
			UserID: userID,
		})

		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusForbidden, fmt.Errorf("you do not own this product"))
			return
		}
	}

	// 3. Handle Database Errors (Connection issues, etc.)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "product deleted successfully",
	})
}
