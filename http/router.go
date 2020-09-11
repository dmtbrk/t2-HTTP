package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/ortymid/t2-http/jwt"
	"github.com/ortymid/t2-http/market"
)

// type ProductsRequest struct{}

type contextKey int

const KeyUserID contextKey = 0

// Router implements standard library http.Handler interface.
// It acts as an entry point to the request handling.
type Router struct {
	Market    market.Interface
	JWTAlg    string
	JWTSecret interface{}
}

// ServeHTTP dispatches incoming http requests to specific handlers.
func (rt *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// All responses are JSON-encoded.
	w.Header().Add("Content-Type", "application/json")

	// Leverage gorilla/mux.
	router := mux.NewRouter()
	rt.registerHandlers(router)

	// Try to get a user id and set it to the request context.
	// No user id means anonymous request.
	req, err := rt.withUserID(req)
	if err != nil {
		err = fmt.Errorf("authorization: %w", err)
		writeError(w, http.StatusForbidden, err)
		return
	}

	// Pass the request to the gorilla/mux router.
	router.ServeHTTP(w, req)
}

func (rt *Router) registerHandlers(r *mux.Router) {
	productHandler := &ProductHandler{
		market: rt.Market,
	}
	s := r.PathPrefix("/products").Subrouter()
	productHandler.RegisterHandlers(s)
}

// withUserID attaches a user ID obtained from the JWT to the request context.
// getToken function defines where is the JWT expected to be found.
func (rt *Router) withUserID(req *http.Request) (*http.Request, error) {
	tokenString, err := getTokenString(req)
	if err != nil {
		return nil, err
	}
	if len(tokenString) == 0 {
		return req, nil // ok, no token
	}

	claims, err := jwt.Parse(tokenString, rt.JWTAlg, rt.JWTSecret)
	if err != nil {
		return nil, fmt.Errorf("parsing request token: %w", err)
	}

	userID := strconv.Itoa(claims.UserID)

	ctx := req.Context()
	ctx = context.WithValue(ctx, KeyUserID, userID)
	return req.WithContext(ctx), nil
}

// getTokenString looks for the JWT in the Authorization header.
// Absence of the token cosidered a normal case.
func getTokenString(req *http.Request) (string, error) {
	auth := req.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", nil // ok, no token
	}

	authFields := strings.Fields(auth)
	if len(authFields) != 2 {
		return "", errors.New("malformed Authorization header")
	}

	typ := authFields[0]
	if !strings.EqualFold(typ, "Bearer") {
		return "", errors.New("Authorization type is not Bearer")
	}

	token := authFields[1]
	return token, nil
}

// writeError writes an error to the response as a JSON-encoded string.
func writeError(w http.ResponseWriter, status int, err error) {
	log.Println("ERROR:", err)

	payload := struct {
		Message string `json:"message"`
	}{
		Message: err.Error(),
	}

	w.WriteHeader(status)
	err = json.NewEncoder(w).Encode(payload)
	if err != nil {
		err = fmt.Errorf("encoding error: %w", err)
		log.Println("ERROR:", err)
	}
}
