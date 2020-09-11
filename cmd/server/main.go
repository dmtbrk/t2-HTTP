package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	httpserver "github.com/ortymid/t2-http/http"
	"github.com/ortymid/t2-http/market"
	httpservice "github.com/ortymid/t2-http/service/http"
	"github.com/ortymid/t2-http/service/mem"
)

type Config struct {
	Port           int
	JWTAlg         string
	JWTSecret      interface{}
	UserServiceURL string
}

func main() {
	config := getConfig()

	userService := httpservice.NewUserService(config.UserServiceURL)
	productService := mem.NewProductService()

	m := &market.Market{
		UserService:    userService,
		ProductService: productService,
	}

	httpserver.Run(config.Port, config.JWTAlg, config.JWTSecret, m)
}

func getConfig() *Config {
	portString := getEnvDefault("PORT", "8080")
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic("cannot read PORT: " + err.Error())
	}

	jwtAlg := getEnvDefault("JWT_ALG", "HS256")
	jwtSecret, err := getKey(os.Getenv("KEY_SERVICE_URL"))
	if err != nil {
		panic(fmt.Errorf("cannot get JWT secret: %w", err))
	}

	usURL := os.Getenv("USER_SERVICE_URL")

	return &Config{
		Port:           port,
		JWTAlg:         jwtAlg,
		JWTSecret:      jwtSecret,
		UserServiceURL: usURL,
	}
}

func getKey(url string) (*rsa.PublicKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("something went wrong")
	}

	key := &rsa.PublicKey{}
	err = json.NewDecoder(resp.Body).Decode(key)
	defer resp.Body.Close()

	return key, err
}

func getEnvDefault(key string, d string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = d
	}
	return val
}
