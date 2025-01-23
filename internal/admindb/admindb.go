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
