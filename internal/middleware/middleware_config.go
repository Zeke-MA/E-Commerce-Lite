package middleware

import (
	"log/slog"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
)

type MiddlewareSiteConfig struct {
	*config.SiteConfig
	*slog.Logger
}
