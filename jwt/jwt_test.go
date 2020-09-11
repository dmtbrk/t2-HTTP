package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"reflect"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

func TestParse(t *testing.T) {
	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	t.Run("Should parse a valid token", func(t *testing.T) {
		wantClaims := &Claims{
			UserID: 1,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, wantClaims)
		tokenString, err := token.SignedString(key)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return
		}

		got, err := Parse(tokenString, "RS256", key.Public())
		if err != nil {
			t.Errorf("Parse() unexpected error: %v", err)
			return
		}
		if !reflect.DeepEqual(got, wantClaims) {
			t.Errorf("Parse() = %v, want %v", got, wantClaims)
		}
	})
}
