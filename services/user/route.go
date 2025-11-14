// This file is user/route.go
package user

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// Routes sets up all user-related API endpoints.
func Routes(database *sql.DB) chi.Router {
	r := chi.NewRouter()
	q := db.New(database)

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, q)
	})
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		handleRegister(w, r, q)
	})

	return r
}

func handleLogin(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	// Parse the request payload
	var payload mytypes.LoginPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	// Get user by email
	user, err := q.GetUserByEmail(context.Background(), payload.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password)); err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with token and user info
	utils.RespondWithJSON(w, http.StatusOK, mytypes.LoginResponse{
		Token: token,
		User: mytypes.UserResponse{
			ID:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	})
}
