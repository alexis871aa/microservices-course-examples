package ufo

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (r *repository) Delete(_ context.Context, uuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	sighting, ok := r.data[uuid]
	if !ok {
		return model.ErrSightingNotFound
	}

	// Мягкое удаление - устанавливаем deleted_at
	sighting.DeletedAt = lo.ToPtr(time.Now())

	r.data[uuid] = sighting

	return nil
}
