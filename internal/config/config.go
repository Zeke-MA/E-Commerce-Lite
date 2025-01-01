package config

import (
	"database/sql"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/database"
)

type SiteConfig struct {
	DbConnection *sql.DB
	DbQueries    *database.Queries
}
