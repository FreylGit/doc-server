package document

import (
	"context"
	"doc-server/internal/models"
	"doc-server/internal/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

func (s *serv) GetList(ctx context.Context, token string, filter map[string]interface{}, limit int64) ([]models.Document, error) {
	// Проверка токена
	claims, err := utils.ParseToken(token, s.configApp.SecretKey())
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	userId, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return nil, err
	}

	// Формирование ключа для кэша
	filterJSON, err := json.Marshal(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal filter: %w", err)
	}
	cacheKey := fmt.Sprintf("user:%d:docs:list:filter:%s:limit:%d", userId, filterJSON, limit)

	cachedData, err := s.cache.Client().Get(ctx, cacheKey).Result()
	if err == nil {
		var documents []models.Document
		if err := json.Unmarshal([]byte(cachedData), &documents); err == nil {
			return documents, nil
		}
	}

	// Получение данных из базы
	list, err := s.docRepo.GetList(ctx, userId, filter, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get documents: %w", err)
	}
	if len(list) == 0 {
		// Получение данных из базы по публичным документам
		list, err = s.docRepo.GetPublicList(ctx, filter, limit)
	}
	// Сохранение результата в кэш
	data, err := json.Marshal(list)
	if err == nil {
		if cacheErr := s.cache.Client().Set(ctx, cacheKey, data, 0).Err(); cacheErr != nil {
			return nil, cacheErr
		}
	}

	return list, nil
}
