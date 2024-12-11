package document

import (
	"context"
	"doc-server/internal/models"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
)

const basePath = "././uploads/"

func (s *serv) Save(ctx context.Context, file multipart.File, document models.Document) (string, error) {
	filePath := basePath + document.Name

	// Сохраняем файл на диск
	destFile, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filePath, err)
		return "", err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		log.Printf("Failed to write file content: %v", err)
		return "", err
	}

	// Сохраняем метаданные в базу
	err = s.docRepo.Create(ctx, document)
	if err != nil {
		log.Printf("Failed to save document metadata: %v", err)
		return "", err
	}

	// Инвалидация кэша списка
	cacheKey := fmt.Sprintf("user:%d:docs:list", document.UserId)
	err = s.cache.Client().Del(ctx, cacheKey).Err()
	if err != nil {
		log.Printf("Failed to invalidate cache: %v", err)
	}

	// Возвращаем содержимое, если это JSON-файл
	switch document.Mime {
	case "json":
		data, err := readFile(document.Name)
		return data, err
	default:
		return "", nil
	}
}
