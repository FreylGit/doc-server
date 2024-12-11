package config

type ConfigHTTP interface {
	Address() string
}

type ConfigPG interface {
	DSN() string
	Settings() PGSettings
}

type ConfigApp interface {
	AdminToken() string
	SecretKey() []byte
}

type ConfigRedis interface {
	Address() string
	Password() string
	DbNum() int
}

type PGSettings struct {
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   int32
	HealthCheckPeriod int32
}
