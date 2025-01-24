package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/admindb"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/utils"
	"github.com/gorilla/mux"
)

type product struct {
	ProductId          string  `json:"product_id"`
	ProductName        string  `json:"product_name"`
	UpcId              string  `json:"upc_id"`
	ProductDescription *string `json:"product_desc"`
	CurrentPrice       string  `json:"price"`
	OnHand             int     `json:"on_hand"`
}

func (cfg *HandlerSiteConfig) AdminAddProduct(w http.ResponseWriter, r *http.Request) {
	type addProduct struct {
		ProductName        string  `json:"product_name"`
		UpcId              string  `json:"upc_id"`
		ProductDescription *string `json:"product_desc"`
		CurrentPrice       string  `json:"price"`
		OnHand             int     `json:"on_hand"`
	}

	requestUserID, ok := utils.GetContextUserID(r.Context())

	if !ok {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	productID := mux.Vars(r)["product_id"]

	if productID == "" {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	authorized, err := admindb.IsUserAdmin(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	product := addProduct{}
	err = decoder.Decode(&product)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	dbAddProductResult, err := admindb.AddProductToSite(productID, product.ProductName, product.UpcId, product.CurrentPrice,
		server.StringToNullString(product.ProductDescription), int32(product.OnHand), requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		if errors.Is(err, admindb.ErrProductAlreadyExists) {
			server.RespondWithError(w, http.StatusConflict, string(server.MsgConflict), err)
			return
		}
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, dbAddProductResult)
}

func (cfg *HandlerSiteConfig) AdminRemoveProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["product_id"]

	if productID == "" {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	requestUserID, ok := utils.GetContextUserID(r.Context())

	if !ok {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	authorized, err := admindb.IsUserAdmin(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	removedProduct, err := admindb.RemoveProductFromSite(productID, r.Context(), cfg.SiteConfig)

	if err != nil {
		if errors.Is(err, admindb.ErrProductNotFound) {
			server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), err)
			return
		}
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, removedProduct)
}

func (cfg *HandlerSiteConfig) AdminChangePrice(w http.ResponseWriter, r *http.Request) {

	type priceChange struct {
		NewPrice string `json:"price"`
	}

	productID := mux.Vars(r)["product_id"]

	if productID == "" {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	requestUserID, ok := utils.GetContextUserID(r.Context())

	if !ok {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	authorized, err := admindb.IsUserAdmin(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	productPriceChange := priceChange{}
	err = decoder.Decode(&productPriceChange)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	dbFoundProduct, err := cfg.DbQueries.FindProduct(r.Context(), productID)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if dbFoundProduct.CurrentPrice == productPriceChange.NewPrice {
		server.RespondWithJSON(w, http.StatusNotModified, nil)
	}

	priceUpdate := database.UpdateProductPriceParams{
		CurrentPrice: productPriceChange.NewPrice,
		ProductID:    productID,
	}

	dbPriceUpdate, err := cfg.DbQueries.UpdateProductPrice(r.Context(), priceUpdate)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	responseUpdate := product{
		ProductId:          productID,
		ProductName:        dbPriceUpdate.ProductName,
		UpcId:              dbPriceUpdate.UpcID,
		ProductDescription: &dbPriceUpdate.ProductDescription.String,
		CurrentPrice:       productPriceChange.NewPrice,
		OnHand:             int(dbPriceUpdate.OnHand),
	}

	server.RespondWithJSON(w, http.StatusOK, responseUpdate)
}

func (cfg *HandlerSiteConfig) AdminAddQuantity(w http.ResponseWriter, r *http.Request) {
	type addQuantity struct {
		Quantity int `json:"add_quantity"`
	}
	productID := mux.Vars(r)["product_id"]

	if productID == "" {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	requestUserID, ok := utils.GetContextUserID(r.Context())

	if !ok {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	authorized, err := admindb.IsUserAdmin(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	productAddQuantity := addQuantity{}
	err = decoder.Decode(&productAddQuantity)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if productAddQuantity.Quantity < 1 {
		server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
		return
	}

	_, err = cfg.DbQueries.FindProduct(r.Context(), productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), err)
			return
		}
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	addQuantityParams := database.UpdateProductQuantityParams{
		OnHand:    int32(productAddQuantity.Quantity),
		ProductID: productID,
	}

	dbUpdateProductQty, err := cfg.DbQueries.UpdateProductQuantity(r.Context(), addQuantityParams)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	addQtyResponse := product{
		ProductId:          productID,
		ProductName:        dbUpdateProductQty.ProductName,
		UpcId:              dbUpdateProductQty.UpcID,
		ProductDescription: &dbUpdateProductQty.ProductDescription.String,
		CurrentPrice:       dbUpdateProductQty.CurrentPrice,
		OnHand:             int(dbUpdateProductQty.OnHand),
	}

	server.RespondWithJSON(w, http.StatusOK, addQtyResponse)

}

func (cfg *HandlerSiteConfig) AdminRemoveQuantity(w http.ResponseWriter, r *http.Request) {
	type removeQuantity struct {
		Quantity int `json:"remove_quantity"`
	}
	productID := mux.Vars(r)["product_id"]

	if productID == "" {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	requestUserID, ok := utils.GetContextUserID(r.Context())

	if !ok {
		server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), nil)
		return
	}

	authorized, err := admindb.IsUserAdmin(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	productRemoveQuantity := removeQuantity{}
	err = decoder.Decode(&productRemoveQuantity)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if productRemoveQuantity.Quantity < 1 {
		server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
		return
	}

	findProduct, err := cfg.DbQueries.FindProduct(r.Context(), productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			server.RespondWithError(w, http.StatusNotFound, string(server.MsgNotFound), err)
			return
		}
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	// if the amount to remove is more than what is in the table send a bad request rather than taking it to 0.
	// Might indicate issue with the data source reporting more on hand than actually available.
	if productRemoveQuantity.Quantity > int(findProduct.OnHand) {
		server.RespondWithError(w, http.StatusBadRequest, string(server.MsgBadRequest), err)
		return
	}

	addQuantityParams := database.UpdateProductQuantityParams{
		OnHand:    int32(productRemoveQuantity.Quantity) * -1,
		ProductID: productID,
	}

	dbUpdateProductQty, err := cfg.DbQueries.UpdateProductQuantity(r.Context(), addQuantityParams)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	addQtyResponse := product{
		ProductId:          productID,
		ProductName:        dbUpdateProductQty.ProductName,
		UpcId:              dbUpdateProductQty.UpcID,
		ProductDescription: &dbUpdateProductQty.ProductDescription.String,
		CurrentPrice:       dbUpdateProductQty.CurrentPrice,
		OnHand:             int(dbUpdateProductQty.OnHand),
	}

	server.RespondWithJSON(w, http.StatusOK, addQtyResponse)
}
