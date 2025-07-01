package ufo

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (r *repository) Update(_ context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, ok := r.data[uuid]
	if !ok {
		return model.ErrSightingNotFound
	}

	// Обновляем поля, только если они были установлены в запросе
	if updateInfo.ObservedAt != nil {
		sighting.Info.ObservedAt = updateInfo.ObservedAt
	}

	if updateInfo.Location != nil {
		sighting.Info.Location = *updateInfo.Location
	}

	if updateInfo.Description != nil {
		sighting.Info.Description = *updateInfo.Description
	}

	if updateInfo.Color != nil {
		sighting.Info.Color = updateInfo.Color
	}

	if updateInfo.Sound != nil {
		sighting.Info.Sound = updateInfo.Sound
	}

	if updateInfo.DurationSeconds != nil {
		sighting.Info.DurationSeconds = updateInfo.DurationSeconds
	}

	sighting.UpdatedAt = lo.ToPtr(time.Now())

	r.data[uuid] = sighting

	return nil
}
