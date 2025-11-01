package user

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	respond "github.com/ARCoder181105/ecom/services"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// Routes sets up all user-related API endpoints.
func Routes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := db.New(database)

	r.Post("/login", handleLogin(q))
	r.Post("/register", handleRegister(q))

	return r
}

// handleLogin handles user login (to be implemented later).
func handleLogin(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respond.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Login successful (to be implemented)",
		})
	}
}

// handleRegister registers a new user in the database.
func handleRegister(q *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
			Username  string `json:"username"`
			Email     string `json:"email"`
			Password  string `json:"password"`
		}

		// Decode request JSON
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			respond.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
			return
		}
		defer r.Body.Close()

		// Basic validation
		if payload.Email == "" || payload.Password == "" || payload.Username == "" {
			respond.RespondWithError(w, http.StatusBadRequest, "Missing required fields")
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
		if err != nil {
			respond.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
			return
		}

		// Create user in database
		user, err := q.CreateUser(context.Background(), db.CreateUserParams{
			FirstName: payload.FirstName,
			LastName:  payload.LastName,
			Username:  payload.Username,
			Email:     payload.Email,
			Password:  string(hashedPassword),
		})
		if err != nil {
			respond.RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		// Success response
		//with Json
		respond.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{
			"message": "User registered successfully",
			"user": mytypes.UserResponse{
				ID:        user.ID.String(),
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Username:  user.Username,
				Email:     user.Email,
				CreatedAt: user.CreatedAt,
			},
		})
	}
}
