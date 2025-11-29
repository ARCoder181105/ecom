package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ARCoder181105/ecom/db"
	"github.com/ARCoder181105/ecom/services/products"
	"github.com/ARCoder181105/ecom/services/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Run() error {

	conn, err := db.NewPostgresStorage()
	if err != nil {
		return fmt.Errorf("❌ failed to connect database: %v", err)
	}
	s.db = conn
	defer s.db.Close()

	log.Println("✅ Database connection established")

	r := chi.NewRouter()

	// CORS Configuration
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL"), "http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/user", user.Routes(s.db))
		api.Mount("/product", products.Routes(s.db))
	})

	// Start server
	log.Printf("🚀 Server running on %s\n", s.addr)
	return http.ListenAndServe(s.addr, r)
}
