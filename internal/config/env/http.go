package env

import (
	"doc-server/internal/config"
	"fmt"
	"log"

	"os"

	"strconv"
)

const (
	HTTP_HOST_ENV = "HTTP_HOST"
	HTTP_PORT_ENV = "HTTP_PORT"
)

type configHTTP struct {
	host string
	port int64
}

func NewConfig() config.ConfigHTTP {
	host := os.Getenv(HTTP_HOST_ENV)
	portStr := os.Getenv(HTTP_PORT_ENV)
	if len(portStr) == 0 || len(host) == 0 {
		log.Fatal("Error parse hhtp config")
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse http port to int64")
	}

	return &configHTTP{
		host: host,
		port: port,
	}
}

func (c configHTTP) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}
