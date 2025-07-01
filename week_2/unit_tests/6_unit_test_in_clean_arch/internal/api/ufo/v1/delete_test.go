package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/model"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

func (s *APISuite) TestDeleteSuccess() {
	var (
		uuid = gofakeit.UUID()

		request = &ufoV1.DeleteRequest{
			Uuid: uuid,
		}
	)

	s.ufoService.On("Delete", s.ctx, uuid).Return(nil)

	res, err := s.api.Delete(s.ctx, request)
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *APISuite) TestDeleteNotFound() {
	var (
		uuid = gofakeit.UUID()

		request = &ufoV1.DeleteRequest{
			Uuid: uuid,
		}
	)

	s.ufoService.On("Delete", s.ctx, uuid).Return(model.ErrSightingNotFound)

	res, err := s.api.Delete(s.ctx, request)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.NotFound, st.Code())
}

func (s *APISuite) TestDeleteServiceError() {
	var (
		serviceErr = gofakeit.Error()
		uuid       = gofakeit.UUID()

		request = &ufoV1.DeleteRequest{
			Uuid: uuid,
		}
	)

	s.ufoService.On("Delete", s.ctx, uuid).Return(serviceErr)

	res, err := s.api.Delete(s.ctx, request)
	s.Require().Error(err)
	s.Require().ErrorIs(err, serviceErr)
	s.Require().Nil(res)
}
