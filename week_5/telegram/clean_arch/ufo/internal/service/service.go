package service

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/model"
)

type UFOService interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
}

type TelegramService interface {
	SendUFONotification(ctx context.Context, uuid string, sighting model.SightingInfo) error
}
