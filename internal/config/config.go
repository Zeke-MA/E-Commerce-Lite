package config

import (
	"database/sql"
)

type siteConfig struct {
	DbConnection *sql.DB
}
