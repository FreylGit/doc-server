package document

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/utils"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

func (s *serv) Delete(ctx context.Context, token string, documentId int64) error {
	claims, err := utils.ParseToken(token, s.configApp.SecretKey())
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	userId, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse userId: %w", err)
	}

	// Удаляем документ из репозитория
	name, err := s.docRepo.Delete(ctx, userId, documentId)
	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}

	// Удаляем файл
	err = removeFile(name)
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	// Удаляем кэшированный документ
	docCacheKey := fmt.Sprintf("user:%d:doc:%d", userId, documentId)
	err = s.cache.Client().Del(ctx, docCacheKey).Err()
	if err != nil {
		return fmt.Errorf("failed to invalidate document cache: %w", err)
	}

	// Удаляем/обновляем списки с документами
	cursor := uint64(0)
	for {
		// Используем SCAN для поиска всех ключей списков
		keys, nextCursor, err := s.cache.Client().Scan(ctx, cursor, fmt.Sprintf("user:%d:docs:list:*", userId), 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan cache keys: %w", err)
		}

		// Обрабатываем найденные ключи
		for _, key := range keys {
			cachedData, err := s.cache.Client().Get(ctx, key).Result()
			if err != nil {
				// Если не удалось получить данные, просто игнорируем
				continue
			}

			var documents []models.Document
			err = json.Unmarshal([]byte(cachedData), &documents)
			if err != nil {
				// Если данные повреждены, игнорируем
				continue
			}

			// Удаляем документ из списка
			updatedDocuments := []models.Document{}
			for _, doc := range documents {
				if doc.Id != documentId {
					updatedDocuments = append(updatedDocuments, doc)
				}
			}

			// Если список изменился, обновляем кэш
			if len(updatedDocuments) != len(documents) {
				if len(updatedDocuments) == 0 {
					// Если список стал пустым, удаляем его из кэша
					_ = s.cache.Client().Del(ctx, key).Err()
				} else {
					// Иначе сохраняем обновлённый список
					data, err := json.Marshal(updatedDocuments)
					if err == nil {
						_ = s.cache.Client().Set(ctx, key, data, 0).Err()
					}
				}
			}
		}

		// Если дошли до конца, выходим
		if nextCursor == 0 {
			break
		}
		cursor = nextCursor
	}

	return nil
}

func removeFile(name string) error {
	filePath := basePath + name
	return os.Remove(filePath)
}
