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

	dbUserExists := database.CheckUsernameEmailUniqueParams{
		Username: params.Username,
		Email:    params.Email,
	}

	exists, err := cfg.DbQueries.CheckUsernameEmailUnique(r.Context(), dbUserExists)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if exists {
		server.RespondWithError(w, http.StatusConflict, string(server.MsgInternalError), err)
		return
	}

	dbParams := database.CreateUserParams{
		Username:       params.Username,
		Email:          params.Email,
		HashedPassword: hashedPassword,
	}

	dbUser, err := cfg.DbQueries.CreateUser(r.Context(), dbParams)

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
