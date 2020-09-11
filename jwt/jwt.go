package jwt

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

func Parse(tokenString string, alg string, key interface{}) (*Claims, error) {
	var claims Claims
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if token.Header["alg"] != alg {
			return nil, errors.New("unexpected token algorithm")
		}
		return key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parsing token: %w", err)
	}
	return &claims, nil
}
