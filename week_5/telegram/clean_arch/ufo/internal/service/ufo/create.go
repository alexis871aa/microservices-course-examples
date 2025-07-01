package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/model"
)

func (s *service) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	uuid, err := s.ufoRepository.Create(ctx, info)
	if err != nil {
		return "", err
	}

	// Отправляем уведомление в Telegram
	if err := s.telegramService.SendUFONotification(ctx, uuid, info); err != nil {
		return "", err
	}

	return uuid, nil
}
