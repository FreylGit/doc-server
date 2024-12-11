package app

import (
	"context"
	"doc-server/internal/api/handlers/auth"
	"doc-server/internal/api/handlers/docs"
	auth3 "doc-server/internal/services/auth"
	document2 "doc-server/internal/services/document"
	"doc-server/internal/storage"
	"doc-server/internal/storage/db/pg/document"
	auth2 "doc-server/internal/storage/db/pg/user"
	redis2 "doc-server/internal/storage/db/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"log"
	"time"

	"doc-server/internal/config"
	"doc-server/internal/config/env"
	"doc-server/internal/services"
)

type service_provider struct {
	configHttp  config.ConfigHTTP
	configPg    config.ConfigPG
	configApp   config.ConfigApp
	configRedis config.ConfigRedis
	db          *pgxpool.Pool
	rdb         *redis.Client
	cache       *redis2.Cache

	authHandler  *auth.AuthHandler
	docsHandler  *docs.DocsHandler
	userRepo     storage.UserRepository
	documentRepo storage.DocumentRepository
	authServ     services.AuthService
	documentServ services.DocumentService
}

func NewServiceProvider() *service_provider {
	return &service_provider{}
}

func (sp *service_provider) ConfigHTTP() config.ConfigHTTP {
	if sp.configHttp == nil {
		sp.configHttp = env.NewConfig()
	}
	return sp.configHttp
}

func (sp *service_provider) ConfigPG() config.ConfigPG {
	if sp.configPg == nil {
		sp.configPg = env.NewConfigPG()
	}
	return sp.configPg
}

func (sp *service_provider) ConfigApp() config.ConfigApp {
	if sp.configApp == nil {
		sp.configApp = env.NewConfigApp()
	}
	return sp.configApp
}

func (sp *service_provider) ConfigRedis() config.ConfigRedis {
	if sp.configRedis == nil {
		sp.configRedis = env.NewConfigRedis()
	}
	return sp.configRedis
}

func (sp *service_provider) DB(ctx context.Context) *pgxpool.Pool {
	if sp.db == nil {
		config, err := pgxpool.ParseConfig(sp.ConfigPG().DSN())
		if err != nil {
			log.Fatalf("Unable to parse connection string: %v", err)
		}
		config.MaxConns = sp.ConfigPG().Settings().MinConns
		config.MinConns = sp.ConfigPG().Settings().MinConns
		config.MaxConnLifetime = time.Minute * time.Duration(sp.ConfigPG().Settings().MaxConnLifetime)
		config.HealthCheckPeriod = time.Minute * time.Duration(sp.ConfigPG().Settings().HealthCheckPeriod)
		db, err := pgxpool.NewWithConfig(ctx, config)

		if err != nil {
			log.Fatalf("Failed to connect to DB: %v", err)
		}
		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("Failed ping DB: %v", err)
		}
		sp.db = db
	}
	return sp.db
}

func (sp *service_provider) RedisClient(ctx context.Context) *redis.Client {
	if sp.rdb == nil {
		rdb := redis.NewClient(&redis.Options{
			Addr:     sp.ConfigRedis().Address(),
			Password: sp.ConfigRedis().Password(),
			DB:       sp.ConfigRedis().DbNum(),
		})

		if err := rdb.Ping(ctx).Err(); err != nil {
			log.Fatalf("failed to connect to Redis: %v", err)
		}
		if err := rdb.FlushDB(ctx).Err(); err != nil {
			log.Fatalf("failed to flush Redis DB: %v", err)
		}
		sp.rdb = rdb
	}

	return sp.rdb
}

func (sp *service_provider) Cache(ctx context.Context) *redis2.Cache {
	if sp.cache == nil {
		sp.cache = redis2.NewCache(sp.RedisClient(ctx))
	}

	return sp.cache
}

func (sp *service_provider) UserRepository(ctx context.Context) storage.UserRepository {
	if sp.userRepo == nil {
		sp.userRepo = auth2.NewUserRepository(sp.DB(ctx))
	}

	return sp.userRepo
}

func (sp *service_provider) DocumentRepository(ctx context.Context) storage.DocumentRepository {
	if sp.documentRepo == nil {
		sp.documentRepo = document.NewDocumentRepository(sp.DB(ctx))
	}

	return sp.documentRepo
}

func (sp *service_provider) Close() {
	if sp.db != nil {
		sp.db.Close()
	}
}

func (sp *service_provider) AuthHandler(ctx context.Context) *auth.AuthHandler {
	if sp.authHandler == nil {
		sp.authHandler = auth.NewSongHandler(sp.AuthService(ctx))
	}

	return sp.authHandler
}

func (sp *service_provider) DocsHandler(ctx context.Context) *docs.DocsHandler {
	if sp.docsHandler == nil {
		sp.docsHandler = docs.NewSongHandler(sp.ConfigApp(), sp.AuthService(ctx), sp.DocumentService(ctx))
	}

	return sp.docsHandler
}

func (sp *service_provider) AuthService(ctx context.Context) services.AuthService {
	if sp.authServ == nil {
		sp.authServ = auth3.NewService(sp.ConfigApp(), sp.UserRepository(ctx))
	}

	return sp.authServ
}

func (sp *service_provider) DocumentService(ctx context.Context) services.DocumentService {
	if sp.documentServ == nil {
		sp.documentServ = document2.NewService(sp.ConfigApp(), sp.UserRepository(ctx), sp.DocumentRepository(ctx), sp.Cache(ctx))
	}

	return sp.documentServ
}
