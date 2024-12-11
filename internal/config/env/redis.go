package env

import (
	cfg "doc-server/internal/config"
	"fmt"
	"log"
	"os"
	"strconv"
)

const (
	REDIS_HOST_ENV     = "REDIS_HOST"
	REDIS_PORT_ENV     = "REDIS_PORT"
	REDIS_DB_ENV       = "REDIS_DB"
	REDIS_PASSWORD_ENV = "REDIS_PASSWORD"
)

type configRedis struct {
	host     string
	port     int64
	db_num   int
	password string
}

func NewConfigRedis() cfg.ConfigRedis {
	host := os.Getenv(REDIS_HOST_ENV)
	portStr := os.Getenv(REDIS_PORT_ENV)
	db_numStr := os.Getenv(REDIS_DB_ENV)
	password := os.Getenv(REDIS_PASSWORD_ENV)
	if len(portStr) == 0 || len(host) == 0 || len(db_numStr) == 0 || len(password) == 0 {
		log.Fatal("Error parse redis config")
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse http port to int64")
	}
	db_num, err := strconv.Atoi(db_numStr)
	if err != nil {
		log.Fatal("Error parse http port to int64")
	}
	return &configRedis{
		host:     host,
		port:     port,
		db_num:   db_num,
		password: password,
	}
}

func (c configRedis) Address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c configRedis) Password() string {
	return c.password
}

func (c configRedis) DbNum() int {
	return c.db_num
}
