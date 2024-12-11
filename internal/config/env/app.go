package env

import (
	"doc-server/internal/config"
	"log"
	"os"
)

const (
	APP_TOKEN_ADMIN_ENV = "APP_TOKEN_ADMIN"
	APP_SECRET_KEY_ENV  = "APP_SECRET_KEY"
)

type configApp struct {
	tokenAdmin string
	secretKey  []byte
}

func NewConfigApp() config.ConfigApp {
	token := os.Getenv(APP_TOKEN_ADMIN_ENV)
	if len(token) == 0 {
		log.Fatalf("Error parse admin token")
	}
	secretKey := []byte(os.Getenv(APP_SECRET_KEY_ENV))
	return &configApp{tokenAdmin: token, secretKey: secretKey}
}

func (c configApp) AdminToken() string {
	return c.tokenAdmin
}

func (c configApp) SecretKey() []byte {
	return c.secretKey
}
