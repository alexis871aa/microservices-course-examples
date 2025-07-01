package repository

import (
	"context"
	"time"

	"github.com/olezhek28/microservices-course-examples/week_6/redis/clean_arch/ufo/internal/model"
)

type UFORepository interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
}

type UFOCacheRepository interface {
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Set(ctx context.Context, uuid string, sighting model.Sighting, ttl time.Duration) error
	Delete(ctx context.Context, uuid string) error
}
