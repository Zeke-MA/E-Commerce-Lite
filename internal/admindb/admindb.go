package admindb

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
	"github.com/google/uuid"
)

var ErrProductAlreadyExists = errors.New("product already exists")
var ErrProductNotFound = errors.New("product not found")
var ErrNoPriceDifference = errors.New("supplied price is the same as the current price")

func IsUserAdmin(userId uuid.UUID, ctx context.Context, cfg *config.SiteConfig) (bool, error) {
	adminCheck, err := cfg.DbQueries.IsAdmin(ctx, userId)

	if err != nil {
		return false, err
	}

	if !adminCheck.IsAdmin {
		return false, nil
	}

	return true, nil

}

func AddProductToSite(productID, productName, upcID, currentPrice string, productDesc sql.NullString,
	onHand int32, createdBy uuid.UUID, ctx context.Context, cfg *config.SiteConfig) (database.Product, error) {

	alreadyExists, err := productExists(productID, ctx, cfg)

	if err != nil {
		return database.Product{}, err
	}

	if alreadyExists {
		return database.Product{}, ErrProductAlreadyExists
	}

	addProductDB := database.AddProductParams{
		ProductID:          productID,
		ProductName:        productName,
		UpcID:              upcID,
		ProductDescription: productDesc,
		CurrentPrice:       currentPrice,
		OnHand:             onHand,
		CreatedBy:          createdBy,
	}

	insertProduct, err := cfg.DbQueries.AddProduct(ctx, addProductDB)

	if err != nil {
		return database.Product{}, err
	}

	return insertProduct, nil
}

func RemoveProductFromSite(productID string, ctx context.Context, cfg *config.SiteConfig) (database.Product, error) {
	productFound, err := productExists(productID, ctx, cfg)

	if err != nil {
		return database.Product{}, err
	}

	if !productFound {
		return database.Product{}, ErrProductNotFound
	}

	removeProduct, err := cfg.DbQueries.RemoveProduct(ctx, productID)

	if err != nil {
		return database.Product{}, err
	}

	return removeProduct, nil

}

func UpdateProductPrice(productID string, newPrice string, ctx context.Context, cfg *config.SiteConfig) (database.Product, error) {
	productFound, err := productDetails(productID, ctx, cfg)

	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return database.Product{}, ErrProductNotFound
		}
		return database.Product{}, err
	}

	if productFound.CurrentPrice == newPrice {
		return database.Product{}, ErrNoPriceDifference
	}

	updateParams := database.UpdateProductPriceParams{
		CurrentPrice: newPrice,
		ProductID:    productID,
	}

	updatedProductPrice, err := cfg.DbQueries.UpdateProductPrice(ctx, updateParams)

	if err != nil {
		return database.Product{}, err
	}

	return updatedProductPrice, nil
}

// Check if product already exists helper to determine if action can or will be taken
func productExists(productID string, ctx context.Context, cfg *config.SiteConfig) (bool, error) {
	_, err := cfg.DbQueries.FindProduct(ctx, productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, ErrProductAlreadyExists
}

func productDetails(productID string, ctx context.Context, cfg *config.SiteConfig) (database.FindProductRow, error) {
	productInformation, err := cfg.DbQueries.FindProduct(ctx, productID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return database.FindProductRow{}, ErrProductNotFound
		}
		return database.FindProductRow{}, err
	}

	return productInformation, nil
}
