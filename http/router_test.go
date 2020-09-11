package http

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/ortymid/t2-http/market"
)

var key *rsa.PrivateKey

func init() {
	key, _ = rsa.GenerateKey(rand.Reader, 2048)
}

type MockMarket struct {
	ProductsRet       []*market.Product
	ProductsErr       error
	ProductRet        *market.Product
	ProductErr        error
	AddProductRet     *market.Product
	AddProductErr     error
	ReplaceProductRet *market.Product
	ReplaceProductErr error
	DeleteProductErr  error
}

func (m MockMarket) Products() ([]*market.Product, error) {
	return m.ProductsRet, m.ProductsErr
}

func (m MockMarket) Product(id int) (*market.Product, error) {
	return m.ProductRet, m.ProductErr
}

func (m MockMarket) AddProduct(p *market.Product, userID string) (*market.Product, error) {
	return m.AddProductRet, m.AddProductErr
}

func (m MockMarket) ReplaceProduct(p *market.Product, userID string) (*market.Product, error) {
	return m.ReplaceProductRet, m.ReplaceProductErr
}

func (m MockMarket) DeleteProduct(id int, userID string) error {
	return m.DeleteProductErr
}

func TestRouter_ServeHTTP(t *testing.T) {
	type fields struct {
		Market    market.Interface
		JWTAlg    string
		JWTSecret interface{}
	}
	tests := []struct {
		name       string
		fields     fields
		req        func() *http.Request
		wantStatus int
		wantBody   []byte
	}{
		{
			name: "Should responde with the products list",
			fields: fields{
				Market: MockMarket{
					ProductsRet: []*market.Product{
						{ID: 1, Name: "p1", Price: 100, Seller: "1"},
					},
				},
				JWTAlg:    "RS256",
				JWTSecret: key.Public(),
			},
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/products/", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte("[{\"id\":1,\"name\":\"p1\",\"price\":100,\"seller\":\"1\"}]\n"),
		},
		{
			name: "Should responde with the product detail",
			fields: fields{
				Market: MockMarket{
					ProductRet: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
				},
				JWTAlg:    "RS256",
				JWTSecret: key.Public(),
			},
			req: func() *http.Request {
				return httptest.NewRequest("GET", "/products/1", nil)
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte("{\"id\":1,\"name\":\"p1\",\"price\":100,\"seller\":\"1\"}\n"),
		},
		{
			name: "Should responde with the new product",
			fields: fields{
				Market: MockMarket{
					AddProductRet: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
				},
				JWTAlg:    "RS256",
				JWTSecret: key.Public(),
			},
			req: func() *http.Request {
				r := httptest.NewRequest("POST", "/products/", strings.NewReader("{\"id\":1,\"name\":\"p1\",\"price\":100}\n"))
				r.Header.Add("Authorization", "Bearer "+testToken(t, 1))
				return r
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte("{\"id\":1,\"name\":\"p1\",\"price\":100,\"seller\":\"1\"}\n"),
		},
		{
			name: "Should responde with the edited product",
			fields: fields{
				Market: MockMarket{
					ReplaceProductRet: &market.Product{ID: 1, Name: "p2", Price: 200, Seller: "1"},
				},
				JWTAlg:    "RS256",
				JWTSecret: key.Public(),
			},
			req: func() *http.Request {
				r := httptest.NewRequest("PUT", "/products/1", strings.NewReader("{\"name\":\"p2\",\"price\":200}\n"))
				r.Header.Add("Authorization", "Bearer "+testToken(t, 1))
				return r
			},
			wantStatus: http.StatusOK,
			wantBody:   []byte("{\"id\":1,\"name\":\"p2\",\"price\":200,\"seller\":\"1\"}\n"),
		},
		{
			name: "Should delete the product",
			fields: fields{
				Market: MockMarket{
					ProductRet: &market.Product{ID: 1, Name: "p1", Price: 100, Seller: "1"},
				},
				JWTAlg:    "RS256",
				JWTSecret: key.Public(),
			},
			req: func() *http.Request {
				r := httptest.NewRequest("DELETE", "/products/1", nil)
				r.Header.Add("Authorization", "Bearer "+testToken(t, 1))
				return r
			},
			wantStatus: http.StatusNoContent,
			wantBody:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Router{
				Market:    tt.fields.Market,
				JWTAlg:    tt.fields.JWTAlg,
				JWTSecret: tt.fields.JWTSecret,
			}

			w := httptest.NewRecorder()
			r := tt.req()
			h.ServeHTTP(w, r)

			resp := w.Result()
			body, _ := ioutil.ReadAll(resp.Body)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("%s Status = %d, want %d", r.URL.Path, resp.StatusCode, tt.wantStatus)
			}
			if !bytes.Equal(body, tt.wantBody) {
				t.Errorf("%s Body = %q, want %q", r.URL.Path, body, tt.wantBody)
			}
		})
	}
}

func testToken(t *testing.T, userID int) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"id": userID})
	tokenString, err := token.SignedString(key)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	return tokenString
}
