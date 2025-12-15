package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hiroBzinga/bun"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/laciferin2024/url-shortner.go/enums"
	"github.com/laciferin2024/url-shortner.go/middlewares"
	"github.com/laciferin2024/url-shortner.go/router"
)

var Module = fx.Options(
	fx.Invoke(
		run, // synchronously

	),
	fx.Provide(initLogrus),
)

type in struct {
	fx.In
	Conf           *viper.Viper
	Middlewares    *middlewares.Middleware
	RouterServices router.Services
	DB             *bun.DB `name:"db"`
}

func run(i in) {
	return
	addr := "0.0.0.0"
	server := &Server{
		i.Middlewares,
		i.RouterServices,
		i.Conf,
		&logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.TextFormatter),
			Hooks:     make(logrus.LevelHooks),
			Level:     logrus.DebugLevel,
		},
		i.DB,
	}

	ginEngine := server.setupRouter()
	server.log.Infoln("running the server on port ", i.Conf.GetString(enums.PORT))

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", addr, i.Conf.GetString(enums.PORT)),
		Handler: ginEngine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			server.log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	server.log.Infoln("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := i.DB.Close(); err != nil {
		server.log.Errorln("db connection failed to close")
	} else {
		server.log.Infoln("db connection closed")
	}

	if err := srv.Shutdown(ctx); err != nil {
		server.log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {

	case <-ctx.Done():
		server.log.Infoln("server exits within 5 seconds.")
	}

	server.log.Infoln("server exited")
	os.Exit(0)
}
