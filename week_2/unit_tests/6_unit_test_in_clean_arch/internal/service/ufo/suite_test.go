package ufo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	ufoRepository *mocks.UFORepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.ufoRepository = mocks.NewUFORepository(s.T())

	s.service = NewService(
		s.ufoRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
