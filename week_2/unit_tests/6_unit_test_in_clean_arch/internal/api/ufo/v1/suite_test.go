package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	ufoService *mocks.UFOService

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.ufoService = mocks.NewUFOService(s.T())

	s.api = NewAPI(
		s.ufoService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(APISuite))
}
