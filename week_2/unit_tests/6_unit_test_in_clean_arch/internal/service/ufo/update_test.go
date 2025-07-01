package ufo

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (s *ServiceSuite) TestUpdateSuccess() {
	var (
		uuid           = gofakeit.UUID()
		newDescription = gofakeit.Paragraph(3, 5, 5, " ")
		newColor       = gofakeit.Color()
		newObservedAt  = time.Now()
		newSound       = false
		newDuration    = int32(120)
		newLocation    = gofakeit.City()

		updateInfo = model.SightingUpdateInfo{
			ObservedAt:      &newObservedAt,
			Location:        &newLocation,
			Description:     &newDescription,
			Color:           &newColor,
			Sound:           &newSound,
			DurationSeconds: &newDuration,
		}
	)

	s.ufoRepository.On("Update", s.ctx, uuid, updateInfo).Return(nil)

	err := s.service.Update(s.ctx, uuid, updateInfo)
	s.NoError(err)
}

func (s *ServiceSuite) TestUpdateRepoError() {
	var (
		repoErr        = gofakeit.Error()
		uuid           = gofakeit.UUID()
		newDescription = gofakeit.Paragraph(3, 5, 5, " ")
		newColor       = gofakeit.Color()

		updateInfo = model.SightingUpdateInfo{
			Description: &newDescription,
			Color:       &newColor,
		}
	)

	s.ufoRepository.On("Update", s.ctx, uuid, updateInfo).Return(repoErr)

	err := s.service.Update(s.ctx, uuid, updateInfo)
	s.Error(err)
	s.ErrorIs(err, repoErr)
}
