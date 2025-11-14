package user

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"golang.org/x/crypto/bcrypt"
)

// 2. Capitalize the function name to export it
func handleRegister(w http.ResponseWriter, r *http.Request, q *db.Queries) {
	// Parse the request payload correctly using &payload
	var payload mytypes.RegisterUserPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	defer r.Body.Close()

	_, err := q.GetUserByEmail(context.Background(), payload.Email)
	if err == nil {
		utils.RespondWithJSON(w, 400, map[string]string{"message": "User already exists"})
		return
	} else if err != sql.ErrNoRows {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.PassWord), 10)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Create user in DB
	user, err := q.CreateUser(context.Background(), db.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
	})
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{
			"message": "Unable to create new user in database",
			"error":   err.Error(),
		})
		return
	}

	// Respond with the created user
	utils.RespondWithJSON(w, http.StatusCreated, mytypes.UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
