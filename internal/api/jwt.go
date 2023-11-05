package api

import (
	"errors"
	"time"

	"github.com/artemiyKew/todo-list-rest-api/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var (
	jwtKey = []byte("secret-key")
)

func generateJWT(u *model.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["sub"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func validateToken(c *fiber.Ctx) (jwt.MapClaims, error) {
	if c.Request().Header.Peek("Token") == nil {
		return nil, errors.New("invalid token")
	}
	tk := string(c.Request().Header.Peek("Token"))
	t, err := jwt.Parse(tk, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot parse")
	}

	if err != nil {
		return nil, err
	}

	return claims, nil
}
