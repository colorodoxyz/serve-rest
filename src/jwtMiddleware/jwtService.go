package jwtMiddleware

/**
 * JWT based authentication package that generates and validates JSON web tokens
 */

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

const (
	secret        = "JWTSECRET"
	tokenLifeSpan = 3
)

/**
 * Validate the json web token provided in authorization header
 */
func ValidateJwt(ctxt *gin.Context) error {
	bearerToken := ctxt.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 && splitToken[0] == "Bearer" {
		tokenString := splitToken[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			return err
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if ok && token.Valid {
			return nil
		}
		return errors.New("Invalid token provided")
	}
	return fmt.Errorf("No JWT provided in authorization header")
}

/**
 * Generate JWT token to match the username passed in (admin)
 */
func GenerateToken(username string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = username
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(tokenLifeSpan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
