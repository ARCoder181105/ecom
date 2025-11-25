package products

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
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

func handleImageUpload(w http.ResponseWriter, r *http.Request) {

	// 1. Parse Multipart Form (Max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("file too big"))
		return
	}

	// 2. Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid file"))
		return
	}
	defer file.Close()

	// 3. Generate a unique filename to prevent overwrites
	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), strings.TrimSuffix(handler.Filename, fileExt))

	// 4. Upload to Cloudinary
	imageUrl, err := utils.UploadToCloudinary(file, uniqueFileName)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// 5. Return the URL
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"image_url": imageUrl})
}

func handleCreateProduct(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	var payload mytypes.CreateProductPayload

	// 1. Parse the JSON body
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	// 2. Validate Price (Convert string to decimal)
	price, err := decimal.NewFromString(payload.Price)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid price format"))
		return
	}

	// 3. Insert into Database
	product, err := q.CreateProduct(context.Background(), database.CreateProductParams{
		Name:        payload.Name,
		Description: payload.Description,
		Image:         sql.NullString{String: payload.Image, Valid: payload.Image != ""},
		Price:         price,
		StockQuantity: int32(payload.StockQuantity),
	})

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// 4. Respond with the new product
	utils.RespondWithJSON(w, http.StatusCreated, product)
}
