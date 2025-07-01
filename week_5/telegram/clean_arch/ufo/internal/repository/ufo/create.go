package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/repository/converter"
	repoModel "github.com/olezhek28/microservices-course-examples/week_5/telegram/clean_arch/ufo/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, info model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()

	sighting := repoModel.Sighting{
		Uuid:      newUUID,
		Info:      repoConverter.SightingInfoToRepoModel(info),
		CreatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, sighting)
	if err != nil {
		return "", err
	}

	return newUUID, nil
}
