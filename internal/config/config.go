package config

import (
	"database/sql"
)

type SiteConfig struct {
	DbConnection *sql.DB
}
