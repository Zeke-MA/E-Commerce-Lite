package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
}

func (cfg *HandlerSiteConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Encapsulate this logic in another function within server
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "error hasshing password: ", err)
		return
	}

	dbParams := database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Username,
		HashedPassword: hashedPassword,
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), dbParams)
	/*
		Add some logic to check if the username and or email already exists
	*/
	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, "error creating user: ", err)
		return
	}

	userResponse := UserResponse{
		ID:        dbUser.ID,
		Email:     dbUser.Email,
		Username:  dbUser.Username,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
	}

	server.RespondWithJSON(w, http.StatusOK, userResponse)

}
