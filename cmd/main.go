package main

import (
	"log"
	"os"

	"github.com/ARCoder181105/ecom/cmd/api"
	"github.com/joho/godotenv"
)

func main() {
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal("⚠️  No .env file found, using system environment")
	}

	port := os.Getenv("PORT")
	

	server := api.NewAPIServer(":"+port)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
