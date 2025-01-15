package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/auth"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/google/uuid"
)

type product struct {
	ProductId          string         `json:"product_id"`
	ProductName        string         `json:"product_name"`
	UpcId              string         `json:"upc_id"`
	ProductDescription sql.NullString `json:"product_desc"`
	CurrentPrice       string         `json:"price"`
	OnHand             int            `json:"on_hand"`
}

func (cfg *HandlerSiteConfig) isUserAdmin(context context.Context, userId uuid.UUID) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(context, userId.String())

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

func (cfg *HandlerSiteConfig) AddProduct(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	requestUserID, err := auth.ValidateJWT(bearerToken, cfg.JWTSecret)

	if err != nil {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	authorized, err := cfg.isUserAdmin(r.Context(), requestUserID)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
	}

	decoder := json.NewDecoder(r.Body)
	product := product{}
	err = decoder.Decode(&product)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	addProduct := database.AddProductParams{
		ProductID:          product.ProductId,
		ProductName:        product.ProductName,
		ProductDescription: product.ProductDescription,
		CurrentPrice:       product.CurrentPrice,
		OnHand:             int32(product.OnHand),
		CreatedBy:          requestUserID,
	}

	insertProduct, err := cfg.DbQueries.AddProduct(r.Context(), addProduct)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	affected, _ := insertProduct.RowsAffected()

	if affected == 0 {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	server.RespondWithJSON(w, http.StatusNoContent, nil)
}
