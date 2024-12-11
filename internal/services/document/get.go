package document

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/utils"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

func (s *serv) Get(ctx context.Context, token string, documentId int64) (string, models.Document, error) {
	claims, err := utils.ParseToken(token, s.configApp.SecretKey())
	if err != nil {
		return "", models.Document{}, fmt.Errorf("invalid token: %w", err)
	}
	cacheKey := fmt.Sprintf("doc:%d", documentId)

	cachedData, err := s.cache.Client().Get(ctx, cacheKey).Result()
	if err == nil {
		var doc models.Document
		if err := json.Unmarshal([]byte(cachedData), &doc); err != nil {
			return "", models.Document{}, fmt.Errorf("failed to unmarshal cached data: %w", err)
		}

		if doc.Mime == "json" {
			file, err := readFile(doc.Name)
			if err != nil {
				_ = s.cache.Client().Del(ctx, cacheKey).Err()
				return "", models.Document{}, fmt.Errorf("failed to read file: %w", err)
			}
			return file, doc, nil
		}

		return "", doc, nil
	} else if err != redis.Nil {
		return "", models.Document{}, fmt.Errorf("redis error: %w", err)
	}

	doc, err := s.docRepo.Get(ctx, claims.Subject, documentId)
	if err != nil {
		return "", models.Document{}, fmt.Errorf("document not found: %w", err)
	}

	docBytes, err := json.Marshal(doc)
	if err == nil {
		_ = s.cache.Client().Set(ctx, cacheKey, docBytes, 0).Err()
	}

	if doc.Mime == "json" {
		file, err := readFile(doc.Name)
		return file, doc, err
	}

	return "", doc, nil
}

func readFile(name string) (string, error) {
	filePath := basePath + name
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
