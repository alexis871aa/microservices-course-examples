package ufo

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		uuid            = gofakeit.UUID()
		location        = gofakeit.City()
		description     = gofakeit.Paragraph(3, 5, 5, " ")
		observedAt      = time.Now()
		color           = gofakeit.Color()
		sound           = true
		durationSeconds = int32(60)
		createdAt       = time.Now()
		updatedAt       = time.Now()

		sighting = model.Sighting{
			Uuid: uuid,
			Info: model.SightingInfo{
				ObservedAt:      &observedAt,
				Location:        location,
				Description:     description,
				Color:           &color,
				Sound:           &sound,
				DurationSeconds: &durationSeconds,
			},
			CreatedAt: createdAt,
			UpdatedAt: &updatedAt,
			DeletedAt: nil,
		}
	)

	s.ufoRepository.On("Get", s.ctx, uuid).Return(sighting, nil)

	res, err := s.service.Get(s.ctx, uuid)
	s.NoError(err)
	s.Equal(sighting, res)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr = gofakeit.Error()
		uuid    = gofakeit.UUID()
	)

	s.ufoRepository.On("Get", s.ctx, uuid).Return(model.Sighting{}, repoErr)

	res, err := s.service.Get(s.ctx, uuid)
	s.Error(err)
	s.ErrorIs(err, repoErr)
	s.Empty(res)
}
