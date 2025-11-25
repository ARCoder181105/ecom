package utils

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
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