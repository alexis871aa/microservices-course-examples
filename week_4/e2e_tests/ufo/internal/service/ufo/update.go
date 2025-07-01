package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_4/e2e_tests/ufo/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	err := s.ufoRepository.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	return nil
}
