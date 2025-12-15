package internal

import (
	"github.com/laciferin2024/url-shortner.go/internal/cache"
	"github.com/laciferin2024/url-shortner.go/internal/cron"
	"github.com/laciferin2024/url-shortner.go/internal/db"
	"github.com/laciferin2024/url-shortner.go/internal/genesis"
	"go.uber.org/fx"
)

var Module = fx.Options(
	cache.Module,
	db.Module,
	genesis.Module,
	cron.Module,
)
