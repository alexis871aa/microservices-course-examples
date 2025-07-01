package v1

import (
	"github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/internal/service"
	ufoV1 "github.com/olezhek28/microservices-course-examples/week_2/6_unit_test_in_clean_arch/pkg/proto/ufo/v1"
)

type api struct {
	ufoV1.UnimplementedUFOServiceServer

	ufoService service.UFOService
}

func NewAPI(ufoService service.UFOService) *api {
	return &api{
		ufoService: ufoService,
	}
}
