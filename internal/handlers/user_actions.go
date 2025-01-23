package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/cartdb"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/server"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/utils"
	"github.com/gorilla/mux"
)

type cartItem struct {
	ProductID    int32     `json:"product_id"`
	Quantity     int32     `json:"quantity"`
	PricePerUnit string    `json:"price_per_unit"`
	CreatedAt    time.Time `json:"added_at"`
}

func (cfg *HandlerSiteConfig) UserAddItemToCart(w http.ResponseWriter, r *http.Request) {
	type addCartItem struct {
		Quantity     int32  `json:"quantity"`
		PricePerUnit string `json:"price_per_unit"`
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

	cartID, err := cartdb.UserCartExistsOrCreate(requestUserID, r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	addItemRequest := addCartItem{}
	err = decoder.Decode(&addItemRequest)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	convertedProductID, err := strconv.ParseInt(productID, 10, 32)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	insertCartItem, err := cartdb.
		UserCartAddItem(cartID, int32(convertedProductID), addItemRequest.Quantity, addItemRequest.PricePerUnit, time.Now().Add(cfg.ItemTimeout), r.Context(), cfg.SiteConfig)

	if err != nil {
		server.RespondWithError(w, http.StatusInternalServerError, string(server.MsgInternalError), err)
		return
	}

	successfulCartItemAdd := cartItem{
		ProductID:    insertCartItem.ProductID,
		Quantity:     insertCartItem.Quantity,
		PricePerUnit: insertCartItem.PricePerUnit,
		CreatedAt:    insertCartItem.CreatedAt.Time,
	}

	server.RespondWithJSON(w, http.StatusOK, successfulCartItemAdd)

}
