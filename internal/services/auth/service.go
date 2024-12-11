package auth

import (
	"doc-server/internal/config"
	"doc-server/internal/services"
	"doc-server/internal/storage"
)

type serv struct {
	configApp config.ConfigApp
	userRepo  storage.UserRepository
}

func NewService(configApp config.ConfigApp, userRepo storage.UserRepository) services.AuthService {
	return &serv{
		configApp: configApp,
		userRepo:  userRepo,
	}
}
