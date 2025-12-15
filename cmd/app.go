package main

import (
	"log"

	"go.uber.org/fx"

	"github.com/laciferin2024/url-shortner.go/config"
	"github.com/laciferin2024/url-shortner.go/handlers"
	"github.com/laciferin2024/url-shortner.go/internal"
	"github.com/laciferin2024/url-shortner.go/middlewares"
	"github.com/laciferin2024/url-shortner.go/router"
	"github.com/laciferin2024/url-shortner.go/server"
	"github.com/laciferin2024/url-shortner.go/services"
)

func serverRun() {
	app := fx.New(
		config.Module,
		internal.Module,
		services.Module,
		middlewares.Module,
		handlers.Module,
		router.Module,
		server.Module,
	)

	app.Run()
	log.Fatal("shttuing down")
}
