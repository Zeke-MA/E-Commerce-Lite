package config

import (
	"database/sql"
	"time"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
)

type SiteConfig struct {
	DbConnection       *sql.DB
	DbQueries          *database.Queries
	RefreshTokenExpiry time.Duration
	JWTExpiry          time.Duration
	JWTSecret          string
}
