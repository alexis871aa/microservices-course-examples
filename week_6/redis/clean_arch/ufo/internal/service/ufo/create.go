package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	// Создаем в MongoDB
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	// Получаем созданную запись для кеширования
	sighting, err := s.ufoRepository.Get(ctx, uuid)
	if err == nil {
		// Кешируем созданную запись (игнорируем ошибки кеширования)
		_ = s.cacheRepository.Set(ctx, uuid, sighting, s.cacheTTL)
	}

	return uuid, nil
}
