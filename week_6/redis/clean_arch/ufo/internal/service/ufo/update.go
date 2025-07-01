package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	// Обновляем в MongoDB
	err := s.ufoRepository.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	// Инвалидируем кеш (игнорируем ошибки)
	_ = s.cacheRepository.Delete(ctx, uuid)

	return nil
}
