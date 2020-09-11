package jwt

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	UserID int `json:"id"`
	jwt.StandardClaims
}
