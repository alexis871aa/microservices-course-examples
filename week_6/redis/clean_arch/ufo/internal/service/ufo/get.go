package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.Sighting, error) {
	// Сначала пытаемся получить из кеша
	sighting, err := s.cacheRepository.Get(ctx, uuid)
	if err == nil {
		return sighting, nil
	}

	// Если нет в кеше или ошибка, идем в MongoDB
	sighting, err = s.ufoRepository.Get(ctx, uuid)
	if err != nil {
		return model.Sighting{}, err
	}

	// Сохраняем в кеш (игнорируем ошибки кеширования)
	_ = s.cacheRepository.Set(ctx, uuid, sighting, s.cacheTTL)

	return sighting, nil
}
