package cartdb

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/google/uuid"
)

func UserCartExistsOrCreate(userID uuid.UUID, ctx context.Context, cfg *config.SiteConfig) (int32, error) {
	userCart, err := cfg.DbQueries.FindCartID(ctx, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			newCart, err := cfg.DbQueries.CreateCart(ctx, userID)
			if err != nil {
				return 0, err
			}
			return newCart.ID, nil
		}
		return 0, err
	}

	return userCart, nil
}

func UserCartAddItem(cartID, productID, quantity int32, pricePerUnit string, itemTimeout time.Time, ctx context.Context, cfg *config.SiteConfig) (database.CartItem, error) {
	var cartItem database.CartItem
	itemFoundInCart, err := cartItemExists(cartID, productID, pricePerUnit, ctx, cfg.DbQueries)

	if err != nil {
		return database.CartItem{}, err
	}

	if !itemFoundInCart {
		cartItem, err = insertCartItem(cartID, productID, quantity, pricePerUnit, itemTimeout, ctx, cfg.DbQueries)

		if err != nil {
			return database.CartItem{}, err
		}
	} else {
		cartItem, err = updateCartItemQuantity(cartID, productID, quantity, pricePerUnit, itemTimeout, ctx, cfg.DbQueries)

		if err != nil {
			return database.CartItem{}, err
		}
	}

	return cartItem, nil
}

func cartItemExists(cartID, productID int32, pricePerUnit string, ctx context.Context, dbQueries *database.Queries) (bool, error) {
	cartItemExists := database.CartItemExistsParams{
		CartID:       cartID,
		ProductID:    productID,
		PricePerUnit: pricePerUnit,
	}

	_, err := dbQueries.CartItemExists(ctx, cartItemExists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

func insertCartItem(cartID, productID, quantity int32, pricePerUnit string, itemTimeout time.Time, ctx context.Context, dbQueries *database.Queries) (database.CartItem, error) {
	addItemToCart := database.AddItemToCartParams{
		CartID:       cartID,
		ProductID:    productID,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		ItemTimeout:  itemTimeout,
	}

	cartItemInsert, err := dbQueries.AddItemToCart(ctx, addItemToCart)

	if err != nil {
		return database.CartItem{}, err
	}

	return cartItemInsert, nil
}

func updateCartItemQuantity(cartID, productID, quantity int32, pricePerUnit string, itemTimeout time.Time, ctx context.Context, dbQueries *database.Queries) (database.CartItem, error) {
	updateCartItemQuantity := database.UpdateCartItemQuantityParams{
		Quantity:     quantity,
		ItemTimeout:  itemTimeout,
		CartID:       cartID,
		ProductID:    productID,
		PricePerUnit: pricePerUnit,
	}

	updateCartItemRow, err := dbQueries.UpdateCartItemQuantity(ctx, updateCartItemQuantity)

	if err != nil {
		return database.CartItem{}, err
	}

	return updateCartItemRow, nil
}
