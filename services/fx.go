package services

import (
	"github.com/laciferin2024/url-shortner.go/services/app"
	"github.com/laciferin2024/url-shortner.go/services/auth"
	"github.com/laciferin2024/url-shortner.go/services/dummy"
	"go.uber.org/fx"
)

var Module = fx.Options(
	dummy.Module,
	app.Module,
	auth.Module,
)
