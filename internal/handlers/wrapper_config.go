package handlers

import (
	"time"

	"github.com/Zeke-MA/E-Commerce-Lite/internal/config"
)

type HandlerSiteConfig struct {
	*config.SiteConfig
	ItemTimeout time.Duration
}
