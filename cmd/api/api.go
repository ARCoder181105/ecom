package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/ARCoder181105/ecom/db"
	"github.com/ARCoder181105/ecom/services/user"
	"github.com/go-chi/chi/v5"
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


	r.Route("/api/v1", func(api chi.Router) {
		api.Mount("/user", user.Routes(s.db)) 
	})

	// Start server
	log.Printf("🚀 Server running on %s\n", s.addr)
	return http.ListenAndServe(s.addr, r)
}
