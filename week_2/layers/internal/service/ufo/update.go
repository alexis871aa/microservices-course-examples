package ufo

import (
	"context"

	"github.com/olezhek28/microservices-course-examples/week_2/layers/internal/model"
)

func (s *service) Update(ctx context.Context, uuid string, updateInfo model.SightingUpdateInfo) error {
	err := s.ufoRepository.Update(ctx, uuid, updateInfo)
	if err != nil {
		return err
	}

	return nil
}
