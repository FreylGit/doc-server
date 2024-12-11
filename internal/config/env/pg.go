package env

import (
	"doc-server/internal/config"
	"fmt"
	"log"

	"os"
	"strconv"
)

const (
	PG_HOST_ENV          = "PG_HOST"
	PG_PORT_ENV          = "PG_PORT"
	PG_USER_ENV          = "PG_USER"
	PG_PASSWORD_ENV      = "PG_PASSWORD"
	PG_DATABASE_NAME_ENV = "PG_DATABASE_NAME"

	PG_MAX_CONNS_ENV           = "PG_MAX_CONNS"
	PG_MIN_CONNS_ENV           = "PG_MIN_CONNS"
	PG_MAX_CONN_LIFETIME_ENV   = "PG_MAX_CONN_LIFETIME"
	PG_HEALTH_CHECK_PERIOD_ENV = "PG_HEALTH_CHECK_PERIOD"
)

type configPG struct {
	host     string
	port     int64
	user     string
	password string
	dbName   string
	settings config.PGSettings
}

func NewConfigPG() config.ConfigPG {
	host := os.Getenv(PG_HOST_ENV)
	portStr := os.Getenv(PG_PORT_ENV)
	user := os.Getenv(PG_USER_ENV)
	pass := os.Getenv(PG_PASSWORD_ENV)
	dbName := os.Getenv(PG_DATABASE_NAME_ENV)

	maxConnsStr := os.Getenv(PG_MAX_CONNS_ENV)
	minConnsStr := os.Getenv(PG_MIN_CONNS_ENV)
	maxConnLifetimeStr := os.Getenv(PG_MAX_CONN_LIFETIME_ENV)
	healthCheckPeriodStr := os.Getenv(PG_HEALTH_CHECK_PERIOD_ENV)
	if len(portStr) == 0 || len(host) == 0 ||
		len(user) == 0 || len(pass) == 0 ||
		len(maxConnsStr) == 0 || len(minConnsStr) == 0 ||
		len(maxConnLifetimeStr) == 0 || len(healthCheckPeriodStr) == 0 {
		log.Fatal("Error parse pg config")
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse pg port to int64")
	}
	maxConns, err := strconv.ParseInt(maxConnsStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse pg max conns to int64")
	}
	minConns, err := strconv.ParseInt(minConnsStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse pg min conns to int64")
	}
	maxConnLifetime, err := strconv.ParseInt(maxConnLifetimeStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse pg max conn lifetime to int64")
	}
	healthCheckPeriod, err := strconv.ParseInt(healthCheckPeriodStr, 10, 64)
	if err != nil {
		log.Fatal("Error parse pg health check period to int64")
	}
	return &configPG{
		host:     host,
		port:     port,
		user:     user,
		password: pass,
		dbName:   dbName,
		settings: config.PGSettings{
			MaxConns:          int32(maxConns),
			MinConns:          int32(minConns),
			MaxConnLifetime:   int32(maxConnLifetime),
			HealthCheckPeriod: int32(healthCheckPeriod),
		},
	}
}

func (c configPG) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", c.host, c.port, c.dbName, c.user, c.password)
}

func (c configPG) Settings() config.PGSettings {
	return c.settings
}
