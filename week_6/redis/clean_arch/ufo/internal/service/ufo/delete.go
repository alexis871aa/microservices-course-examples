package ufo

import (
	"context"
)

func (s *service) Delete(ctx context.Context, uuid string) error {
	// Удаляем в MongoDB (мягкое удаление)
	err := s.ufoRepository.Delete(ctx, uuid)
	if err != nil {
		return err
	}

	// Инвалидируем кеш (игнорируем ошибки)
	_ = s.cacheRepository.Delete(ctx, uuid)

	return nil
}
