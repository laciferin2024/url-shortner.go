package router

import (
	"github.com/laciferin2024/url-shortner.go/handlers"
	"github.com/laciferin2024/url-shortner.go/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

//type services struct {
//	Root NoAuthRouting
//}

type service struct {
	conf         *viper.Viper
	middlewares  *middlewares.Middleware
	homeHandler  *handlers.HomeHandler
	dummyHandler *handlers.DummyHandler
	appHandler   *handlers.AppHandler
}

type Middleware func() gin.HandlerFunc
