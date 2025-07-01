package repository

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_7/tracing/ufo/internal/model"
)

type UFORepository interface {
	Create(ctx context.Context, info model.SightingInfo) (string, error)
	Get(ctx context.Context, uuid string) (model.Sighting, error)
	Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error
	Delete(ctx context.Context, uuid string) error
}
