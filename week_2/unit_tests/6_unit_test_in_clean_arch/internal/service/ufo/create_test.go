package ufo

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
)

func (s *ServiceSuite) TestCreateSuccess() {
	var (
		location        = gofakeit.City()
		description     = gofakeit.Paragraph(3, 5, 5, " ")
		observedAt      = time.Now()
		color           = gofakeit.Color()
		sound           = true
		durationSeconds = int32(60)

		expectedUUID = gofakeit.UUID()

		sightingInfo = model.SightingInfo{
			ObservedAt:      &observedAt,
			Location:        location,
			Description:     description,
			Color:           &color,
			Sound:           &sound,
			DurationSeconds: &durationSeconds,
		}
	)

	s.ufoRepository.On("Create", s.ctx, sightingInfo).Return(expectedUUID, nil)

	uuid, err := s.service.Create(s.ctx, sightingInfo)
	s.Require().NoError(err)
	s.Require().Equal(expectedUUID, uuid)
}

func (s *ServiceSuite) TestCreateRepoError() {
	var (
		repoErr         = gofakeit.Error()
		location        = gofakeit.City()
		description     = gofakeit.Paragraph(3, 5, 5, " ")
		observedAt      = time.Now()
		color           = gofakeit.Color()
		sound           = true
		durationSeconds = int32(60)

		sightingInfo = model.SightingInfo{
			ObservedAt:      &observedAt,
			Location:        location,
			Description:     description,
			Color:           &color,
			Sound:           &sound,
			DurationSeconds: &durationSeconds,
		}
	)

	s.ufoRepository.On("Create", s.ctx, sightingInfo).Return("", repoErr)

	uuid, err := s.service.Create(s.ctx, sightingInfo)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(uuid)
}
