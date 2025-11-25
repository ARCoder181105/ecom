package utils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

func UploadToCloudinary(file interface{}, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Initialize Cloudinary using the CLOUDINARY_URL from .env
	cld, err := cloudinary.NewFromURL(os.Getenv("CLOUDINARY_URL"))
	if err != nil {
		return "", fmt.Errorf("failed to intialize cloudinary: %v", err)
	}

	// Upload the file
	// We use "ecom_products" as the folder name in your Cloudinary media library
	uploadParams := uploader.UploadParams{
		Folder:   "ecom_products",
		PublicID: filename,
	}

	result, err := cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	return result.SecureURL, nil
}

func HandleImageUpload(w http.ResponseWriter, r *http.Request) {

	// 1. Parse Multipart Form (Max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("file too big"))
		return
	}

	// 2. Retrieve the file
	file, handler, err := r.FormFile("image")
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Errorf("invalid file"))
		return
	}
	defer file.Close()

	// 3. Generate a unique filename to prevent overwrites
	fileExt := filepath.Ext(handler.Filename)
	uniqueFileName := fmt.Sprintf("%s-%s", uuid.New().String(), strings.TrimSuffix(handler.Filename, fileExt))

	// 4. Upload to Cloudinary
	imageUrl, err := UploadToCloudinary(file, uniqueFileName)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// 5. Return the URL
	RespondWithJSON(w, http.StatusOK, map[string]string{"image_url": imageUrl})
}
