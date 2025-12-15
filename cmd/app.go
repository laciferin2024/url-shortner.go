package main

import (
	"log"

	"go.uber.org/fx"

	"github.com/goferHiro/url-shortner/config"
	"github.com/goferHiro/url-shortner/handlers"
	"github.com/goferHiro/url-shortner/internal"
	"github.com/goferHiro/url-shortner/middlewares"
	"github.com/goferHiro/url-shortner/router"
	"github.com/goferHiro/url-shortner/server"
	"github.com/goferHiro/url-shortner/services"
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
