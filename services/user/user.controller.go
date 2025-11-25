package user

import (
	"context"
	"database/sql"
	"net/http"

	database "github.com/ARCoder181105/ecom/db/migrate/sqlc"
	mytypes "github.com/ARCoder181105/ecom/types"
	"github.com/ARCoder181105/ecom/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Internal function to generate and set access token
func setAccessTokenCookie(w http.ResponseWriter, userID, email string, role string) (string, error) {
	token, err := utils.GenerateJWT(userID, email, role)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "accessToken",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   3600 * 24,
	})

	return token, nil
}

func handleRegister(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	// Parse the request payload correctly using &payload
	var payload mytypes.RegisterUserPayload
	if err := utils.ParseJson(r, &payload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

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
	user, err := q.CreateUser(context.Background(), database.CreateUserParams{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  string(hashedPassword),
		Role:      database.UserRoleCustomer,
	})
	if err != nil {
		utils.RespondWithJSON(w, 400, map[string]string{
			"message": "Unable to create new user in database",
			"error":   err.Error(),
		})
		return
	}

	token, err := setAccessTokenCookie(w, user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with the created user
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
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

func handleLogin(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	var loginUserPayload mytypes.LoginUserPayload

	if err := utils.ParseJson(r, &loginUserPayload); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	user, err := q.GetUserByEmail(context.Background(), loginUserPayload.Email)
	if err == sql.ErrNoRows {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		return
	} else if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginUserPayload.PassWord))
	if err != nil {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid credentials"})
		return
	}

	// Generate and set JWT token
	token, err := setAccessTokenCookie(w, user.ID.String(), user.Email, string(user.Role))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Respond with token and user info
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
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

func handleProfile(w http.ResponseWriter, r *http.Request, q *database.Queries) {
	claims, err := utils.GetClaims(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err)
		return
	}

	userId := claims.UserID

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	user, err := q.GetUserByID(context.Background(), userUUID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, mytypes.UserResponse{
		ID:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	})
}
