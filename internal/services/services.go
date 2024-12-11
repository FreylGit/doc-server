package services

import (
	"context"
	"doc-server/internal/models"
	"mime/multipart"
)

type AuthService interface {
	Register(ctx context.Context, token string, login string, pswd string) error
	Auth(ctx context.Context, login string, pswd string) (string, error)
	SignOut(ctx context.Context, token string)
}

type DocumentService interface {
	Save(ctx context.Context, file multipart.File, document models.Document) (string, error)
	GetList(ctx context.Context, token string, filter map[string]interface{}, limit int64) ([]models.Document, error)
	Get(ctx context.Context, token string, documentId int64) (string, models.Document, error)
	Delete(ctx context.Context, token string, documentId int64) error
}
