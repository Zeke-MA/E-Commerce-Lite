package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

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

	authorized, err := cfg.IsUserAdmin(r.Context(), requestUserID)

	if err != nil {
		log.Print("failed admin check")
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
	}

	dbFoundProduct, err := cfg.DbQueries.FindProduct(r.Context(), productID)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		log.Print("check prod")
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if dbFoundProduct.ProductID == productID {
		server.RespondWithError(w, http.StatusConflict, string(server.MsgConflict), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	product := addProduct{}
	err = decoder.Decode(&product)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	dbProduct := database.AddProductParams{
		ProductID:          productID,
		ProductName:        product.ProductName,
		UpcID:              product.UpcId,
		ProductDescription: server.StringToNullString(product.ProductDescription),
		CurrentPrice:       product.CurrentPrice,
		OnHand:             int32(product.OnHand),
		CreatedBy:          requestUserID,
	}

	insertProduct, err := cfg.DbQueries.AddProduct(r.Context(), dbProduct)

	if err != nil {
		log.Print("check add prod")
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	affected, _ := insertProduct.RowsAffected()

	if affected == 0 {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	server.RespondWithJSON(w, http.StatusOK, dbProduct)
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

	authorized, err := cfg.IsUserAdmin(r.Context(), requestUserID)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	if !authorized {
		server.RespondWithError(w, http.StatusUnauthorized, string(server.MsgUnauthorized), err)
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

	removedProduct, err := cfg.DbQueries.RemoveProduct(r.Context(), productID)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	response := product{
		ProductId:          removedProduct.ProductID,
		ProductName:        removedProduct.ProductName,
		UpcId:              removedProduct.UpcID,
		ProductDescription: &removedProduct.ProductDescription.String,
		CurrentPrice:       removedProduct.CurrentPrice,
		OnHand:             int(removedProduct.OnHand),
	}

	server.RespondWithJSON(w, http.StatusOK, response)
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

	authorized, err := cfg.IsUserAdmin(r.Context(), requestUserID)

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

	authorized, err := cfg.IsUserAdmin(r.Context(), requestUserID)

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

	authorized, err := cfg.IsUserAdmin(r.Context(), requestUserID)

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
