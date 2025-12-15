package router

import (
	"github.com/gin-gonic/gin"
	"github.com/laciferin2024/url-shortner.go/middlewares"
)

func (s *service) RoutesWithNoAuth(r *gin.RouterGroup, mws ...Middleware) {

	for _, i := range mws {
		r.Use(i())
	}

	r.GET("/", s.homeHandler.Home)

	r.POST("/dummy", s.dummyHandler.Dummy)

	appRouter := r.Group("/app")

	appRouter.GET("/", s.appHandler.Home)

	appRouter.POST("/url", middlewares.RateLimitMiddleware(), s.appHandler.ShortenUrl)
	appRouter.GET("/url/:surl", s.appHandler.RetrieveUrl)

	r.GET("/admin/urls", s.appHandler.ListUrls)

}
