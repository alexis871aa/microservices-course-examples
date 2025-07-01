package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/converter"
	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (s *APISuite) TestGetSuccess() {
	var (
		uuid            = gofakeit.UUID()
		location        = gofakeit.City()
		description     = gofakeit.Paragraph(3, 5, 5, " ")
		observedAt      = time.Now()
		color           = gofakeit.Color()
		sound           = true
		durationSeconds = int32(60)
		createdAt       = time.Now()

		req = &ufoV1.GetRequest{
			Uuid: uuid,
		}

		modelSighting = model.Sighting{
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
		}

		expectedProtoSighting = converter.SightingToProto(modelSighting)
	)

	s.ufoService.On("Get", s.ctx, uuid).Return(modelSighting, nil)

	res, err := s.api.Get(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(expectedProtoSighting.Uuid, res.GetSighting().GetUuid())
	s.Require().Equal(expectedProtoSighting.Info.Location, res.GetSighting().GetInfo().GetLocation())
	s.Require().Equal(expectedProtoSighting.Info.Description, res.GetSighting().GetInfo().GetDescription())
}

func (s *APISuite) TestGetNotFound() {
	var (
		uuid = gofakeit.UUID()

		req = &ufoV1.GetRequest{
			Uuid: uuid,
		}
	)

	s.ufoService.On("Get", s.ctx, uuid).Return(model.Sighting{}, model.ErrSightingNotFound)

	res, err := s.api.Get(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestGetServiceError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		req = &ufoV1.GetRequest{
			Uuid: uuid,
		}
	)

	s.ufoService.On("Get", s.ctx, uuid).Return(model.Sighting{}, serviceErr)

	res, err := s.api.Get(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
