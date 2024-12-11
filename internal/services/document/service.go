package document

import (
	"doc-server/internal/config"
	"doc-server/internal/services"
	"doc-server/internal/storage"
	"doc-server/internal/storage/db/redis"
)

type serv struct {
	configApp config.ConfigApp
	userRepo  storage.UserRepository
	docRepo   storage.DocumentRepository
	cache     *redis.Cache
}

func NewService(configApp config.ConfigApp,
	userRepo storage.UserRepository,
	docRepo storage.DocumentRepository,
	cache *redis.Cache,
) services.DocumentService {
	return &serv{
		configApp: configApp,
		userRepo:  userRepo,
		docRepo:   docRepo,
		cache:     cache,
	}
}
