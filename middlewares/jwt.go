package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	models_enums "github.com/laciferin2024/url-shortner.go/models/models-enums"
	log "github.com/sirupsen/logrus"
)

func (m *Middleware) AuthorizeJWT() gin.HandlerFunc {

	return func(c *gin.Context) {

		//goland:noinspection GoSnakeCaseUsage
		const BEARER_SCHEMA = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]

		token, err := m.authServices.ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			c.Set(models_enums.AppUserID, claims["ID"])

			m.log.Info("Claims:", claims)
			c.Next()
		} else {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
