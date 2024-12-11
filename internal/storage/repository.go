package storage

import (
	"context"
	"doc-server/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	Get(ctx context.Context, login string) (models.User, error)
	GetById(ctx context.Context, id int64) (models.User, error)
}

type DocumentRepository interface {
	Create(ctx context.Context, document models.Document) error
	GetList(ctx context.Context, userId int64, filter map[string]interface{}, limit int64) ([]models.Document, error)
	GetPublicList(ctx context.Context, filter map[string]interface{}, limit int64) ([]models.Document, error)
	Get(ctx context.Context, login string, documentId int64) (models.Document, error)
	Delete(ctx context.Context, userId int64, documentId int64) (string, error)
}
