package v1

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/converter"
	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (s *APISuite) TestUpdateSuccess() {
	var (
		uuid           = gofakeit.UUID()
		newDescription = gofakeit.Paragraph(3, 5, 5, " ")
		newColor       = gofakeit.Color()
		newObservedAt  = time.Now()

		protoUpdateInfo = &ufoV1.SightingUpdateInfo{
			ObservedAt:  timestamppb.New(newObservedAt),
			Description: wrapperspb.String(newDescription),
			Color:       wrapperspb.String(newColor),
		}

		req = &ufoV1.UpdateRequest{
			Uuid:       uuid,
			UpdateInfo: protoUpdateInfo,
		}

		expectedModelUpdateInfo = converter.UpdateInfoToModel(protoUpdateInfo)
	)

	s.ufoService.On("Update", s.ctx, uuid, expectedModelUpdateInfo).Return(nil)

	res, err := s.api.Update(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *APISuite) TestUpdateNilInfo() {
	var (
		uuid = gofakeit.UUID()

		req = &ufoV1.UpdateRequest{
			Uuid:       uuid,
			UpdateInfo: nil,
		}
	)

	res, err := s.api.Update(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestUpdateNotFound() {
	var (
		uuid           = gofakeit.UUID()
		newDescription = gofakeit.Paragraph(3, 5, 5, " ")
		newColor       = gofakeit.Color()

		protoUpdateInfo = &ufoV1.SightingUpdateInfo{
			Description: wrapperspb.String(newDescription),
			Color:       wrapperspb.String(newColor),
		}

		req = &ufoV1.UpdateRequest{
			Uuid:       uuid,
			UpdateInfo: protoUpdateInfo,
		}

		expectedModelUpdateInfo = converter.UpdateInfoToModel(protoUpdateInfo)
	)

	s.ufoService.On("Update", s.ctx, uuid, expectedModelUpdateInfo).Return(model.ErrSightingNotFound)

	res, err := s.api.Update(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestUpdateServiceError() {
	var (
		serviceErr     = gofakeit.Error()
		uuid           = gofakeit.UUID()
		newDescription = gofakeit.Paragraph(3, 5, 5, " ")
		newColor       = gofakeit.Color()

		protoUpdateInfo = &ufoV1.SightingUpdateInfo{
			Description: wrapperspb.String(newDescription),
			Color:       wrapperspb.String(newColor),
		}

		req = &ufoV1.UpdateRequest{
			Uuid:       uuid,
			UpdateInfo: protoUpdateInfo,
		}

		expectedModelUpdateInfo = converter.UpdateInfoToModel(protoUpdateInfo)
	)

	s.ufoService.On("Update", s.ctx, uuid, expectedModelUpdateInfo).Return(serviceErr)

	res, err := s.api.Update(s.ctx, req)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
