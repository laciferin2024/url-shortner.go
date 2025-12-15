package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/laciferin2024/url-shortner.go/enums"
)

type Services interface {
	GenerateToken(name string, admin bool) string
	ValidateToken(tokenString string) (*jwt.Token, error)
}

func (s *service) GenerateToken(username string, admin bool) (tokenStr string) {

	var err error

	issuer := s.Conf.GetString(enums.JWT_ISSUER)
	secret := s.Conf.GetString(enums.JWT_SECRET)
	expiryInterval := s.Conf.GetDuration(enums.JWT_EXPIRY_INTERVAL)

	now := time.Now().UTC()

	claims := &jwtCustomClaims{
		Issuer:    issuer,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(expiryInterval).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err = token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return
}

func (s *service) ValidateToken(tokenStr string) (token *jwt.Token, err error) {

	secret := s.Conf.GetString(enums.JWT_SECRET)

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}
