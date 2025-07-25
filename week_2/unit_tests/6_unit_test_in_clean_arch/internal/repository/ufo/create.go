package ufo

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
	repoConverter "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository/converter"
	repoModel "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository/model"
)

func (r *repository) Create(_ context.Context, info model.SightingInfo) (string, error) {
	newUUID := uuid.NewString()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[newUUID] = repoModel.Sighting{
		Uuid:      newUUID,
		Info:      repoConverter.SightingInfoToRepoModel(info),
		CreatedAt: time.Now(),
	}

	return newUUID, nil
}
